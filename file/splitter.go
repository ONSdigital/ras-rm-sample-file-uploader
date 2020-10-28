package file

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	logger "logging"
	"mime/multipart"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"cloud.google.com/go/pubsub"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
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
	ciCount, totalUnits, buf := readFileForCountTotals(file)
	sampleSummary, err := f.getSampleSummary(ciCount, totalUnits)
	if err != nil {
		return nil, err
	}
	logger.Info("File uploaded",
		zap.String("filename", handler.Filename),
		zap.Int64("filesize", handler.Size),
		zap.Any("MIMEHeader", handler.Header))
	errorCount := f.Publish(bufio.NewScanner(buf))
	if errorCount > 0 {
		return nil, errors.New("unable to process all of sample file")
	}
	return sampleSummary, nil
}

func (f *FileProcessor) Publish(scanner *bufio.Scanner) int {
	logger.Info("about to publish message",
		zap.String("topic", f.Config.Pubsub.TopicId),
		zap.String("project", f.Config.Pubsub.ProjectId),
	)
	topic := f.Client.Topic(f.Config.Pubsub.TopicId)
	var errorCount = 0
	var wg sync.WaitGroup
	var mux sync.Mutex
	for scanner.Scan() {
		line := scanner.Text()
		logger.Debug("Publishing csv line",
			zap.String("line", line))

		wg.Add(1)
		go func(line string, topic *pubsub.Topic, wg *sync.WaitGroup, mux *sync.Mutex, errorCount *int) {
			defer wg.Done()

			id, err := topic.Publish(f.Ctx, &pubsub.Message{
				Data: []byte(line),
				Attributes: map[string]string{
					"sample_summary_id": f.SampleSummary.Id,
				},
			}).Get(f.Ctx)
			if err != nil {
				logger.Error("Error publishing csv line",
					zap.Error(err),
					zap.String("line", line))
				mux.Lock()
				*errorCount++
				mux.Unlock()
			}
			logger.Debug("csv line delivered",
				zap.String("line", line),
				zap.String("messageId", id))
		}(line, topic, &wg, &mux, &errorCount)
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		logger.Error("Error scanning file",
			zap.Error(err))
	}
	return errorCount
}

func (f *FileProcessor) getSampleSummary(ciCount int, totalUnits int) (*SampleSummary, error) {
	baseUrl := f.Config.Sample.BaseUrl
	logger.Info("about to create sample", zap.String("url", baseUrl+"/samples/samplesummary"))
	summaryRequest := &SampleSummary{
		TotalSampleUnits:              totalUnits,
		ExpectedCollectionInstruments: ciCount,
	}

	b, err := json.Marshal(summaryRequest)
	if err != nil {
		logger.Error("Error marshalling Sample Summary Request",
			zap.Any("summaryRequest", summaryRequest),
			zap.Error(err))
		return nil, err
	}

	resp, err := http.Post(baseUrl+"/samples/samplesummary", "application/json", bytes.NewReader(b))
	if err != nil {
		logger.Error("unable to create a sample summary",
			zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logger.Info("returned sample summary data",
		zap.String("body", string(body)))
	sampleSummary := &SampleSummary{}
	err = json.Unmarshal(body, sampleSummary)
	if err != nil {
		logger.Error("error marshalling response data",
			zap.Error(err))
		return nil, err
	}
	logger.Info("created sample summary",
		zap.Any("samplesummary", sampleSummary))
	f.SampleSummary = sampleSummary
	return sampleSummary, nil
}
