package file

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	log "github.com/sirupsen/logrus"
)

type FileProcessor struct {
	Config        config.Config
	Client        *pubsub.Client
	Ctx           context.Context
	SampleSummary *SampleSummary
}

type SampleSummary struct {
	Id                            string `json:"id"`
	TotalSampleUnits              int    `json:"totalSampleUnits"`
	ExpectedCollectionInstruments int    `json:"expectedCollectionInstruments"`
}

func (f *FileProcessor) ChunkCsv(file multipart.File, handler *multipart.FileHeader) (*SampleSummary, error) {
	sampleSummary, err := f.getSampleSummary()
	if err != nil {
		return nil, err
	}
	log.WithField("filename", handler.Filename).
		WithField("filesize", handler.Size).
		WithField("MIMEHeader", handler.Header).
		Info("File uploaded")
	lineCount, errorCount := f.Publish(bufio.NewScanner(file))
	if errorCount > 0 {
		return nil, errors.New("unable to process all of sample file")
	}
	f.updateSampleSummary(1, lineCount, sampleSummary.Id)
	return sampleSummary, nil
}

func (f *FileProcessor) getCount(scanner *bufio.Scanner) (int, int) {
	sampleCount := 0
	formTypes := make(map[string]string)
	for scanner.Scan() {
		sampleCount++
		line := scanner.Text()
		s := strings.Split(line, ":")
		formTypes[s[26]] = s[26]
	}
	return len(formTypes), sampleCount
}

func (f *FileProcessor) Publish(scanner *bufio.Scanner) (int, int) {
	log.WithField("topic", f.Config.Pubsub.TopicId).
		WithField("project", f.Config.Pubsub.ProjectId).
		Info("about to publish message")
	topic := f.Client.Topic(f.Config.Pubsub.TopicId)
	var errorCount = 0
	var lineCount = 0
	var wg sync.WaitGroup
	var mux sync.Mutex
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		log.WithField("line", line).
			Debug("Publishing csv line")

		wg.Add(1)
		go func(line string, topic *pubsub.Topic, wg *sync.WaitGroup, mux *sync.Mutex) {
			defer wg.Done()

			id, err := topic.Publish(f.Ctx, &pubsub.Message{
				Data: []byte(line),
				Attributes: map[string]string{
					"sample_summary_id": f.SampleSummary.Id,
				},
			}).Get(f.Ctx)
			if err != nil {
				log.WithField("line", line).
					WithError(err).
					Error("Error publishing csv line")
				mux.Lock()
				errorCount++
				mux.Unlock()
			}
			log.WithField("line", line).
				WithField("messageId", id).
				Debug("csv line acknowledged")
		}(line, topic, &wg, &mux)
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.WithError(err).Error("Error scanning file")
	}
	return lineCount, errorCount
}

func (f *FileProcessor) getSampleSummary() (*SampleSummary, error) {
	baseUrl := f.Config.Sample.BaseUrl
	log.WithField("url", baseUrl+"/samples/samplesummary").Info("about to create sample")

	resp, err := http.Post(baseUrl+"/samples/samplesummary", "application/json", nil)
	if err != nil {
		log.WithError(err).Error("unable to create a sample summary")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.WithField("body", string(body)).Info("returned sample summary data")
	sampleSummary := &SampleSummary{}
	err = json.Unmarshal(body, sampleSummary)
	if err != nil {
		log.WithError(err).Error("error marshalling response data")
		return nil, err
	}
	log.WithField("samplesummary", sampleSummary).Info("created sample summary")
	f.SampleSummary = sampleSummary
	return sampleSummary, nil
}

func (f *FileProcessor) updateSampleSummary(ciCount int, totalUnits int, summaryId string) (*SampleSummary, error) {
	baseUrl := f.Config.Sample.BaseUrl
	log.WithField("SampleSummary", summaryId).Info("Updating Sample")
	summaryRequest := &SampleSummary{
		TotalSampleUnits:              totalUnits,
		ExpectedCollectionInstruments: ciCount,
	}

	b, err := json.Marshal(summaryRequest)
	if err != nil {
		log.WithField("summaryRequest", summaryRequest).WithError(err).Error("Error marshalling Sample Summary Request")
		return nil, err
	}

	resp, err := http.NewRequest("PUT", baseUrl+"/samplesummary/"+summaryId, bytes.NewReader(b))
	if err != nil {
		log.WithError(err).Error("unable to update a sample summary")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.WithField("body", string(body)).Info("returned sample summary data")
	sampleSummary := &SampleSummary{}
	err = json.Unmarshal(body, sampleSummary)
	if err != nil {
		log.WithError(err).Error("error marshalling response data")
		return nil, err
	}
	log.WithField("samplesummary", sampleSummary).Info("created sample summary")
	f.SampleSummary = sampleSummary
	return sampleSummary, nil
}
