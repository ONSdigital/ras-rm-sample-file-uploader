package file

import (
	"bufio"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"net/http"
	"sync"
)

type FileProcessor struct {
	Config config.Config
	Client *pubsub.Client
	Ctx context.Context
	SampleSummary string
}

type SampleSummary struct {
	Id string
}

func (f *FileProcessor) ChunkCsv(file multipart.File, handler *multipart.FileHeader) {
	f.getSampleSummary()
	log.WithField("filename", handler.Filename).
		WithField("filesize", handler.Size).
		WithField("MIMEHeader", handler.Header).
		Info("File uploaded")
	f.Publish(bufio.NewScanner(file))
}

func (f *FileProcessor) Publish(scanner *bufio.Scanner) int {
	topic := f.Client.Topic(f.Config.Pubsub.TopicId)
	var errorCount = 0
	var wg sync.WaitGroup
	var mux sync.Mutex
	for scanner.Scan() {
		line := scanner.Text()
		log.WithField("line", line).
			Debug("Publishing csv line")

		wg.Add(1)
		go func(line string, topic *pubsub.Topic, wg *sync.WaitGroup, mux *sync.Mutex, errorCount *int) {
			defer wg.Done()

			id, err := topic.Publish(f.Ctx, &pubsub.Message{
				Data: []byte(line),
				Attributes: map[string]string{
					"sampleSummary":   f.SampleSummary.Id,
				},
			}).Get(f.Ctx)
			if err != nil {
				log.WithField("line", line).
					WithError(err).
					Error("Error publishing csv line")
				mux.Lock()
				*errorCount++
				mux.Unlock()
			}
			log.WithField("line", line).
				WithField("messageId", id).
				Debug("csv line acknowledged")
		}(line, topic, &wg, &mux, &errorCount)
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.WithError(err).Error("Error scanning file")
	}
	return errorCount
}


func (f *FileProcessor) getSampleSummary() string {
	resp, err := http.Post(f.Config.Sample.BaseUrl + "/samples/samplesummary", "\"application/json", nil)
	//resp, err := http.Post("/samples/samplesummary", "\"application/json", nil)
	if err != nil {
		log.WithError(err).Error("Unable to create a sample summary")
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.WithField("body", string(body)).Info("Returned sample summary")
	sampleSummary := &SampleSummary{}
	json.Unmarshal(body, sampleSummary)

	return sampleSummary.Id
}
