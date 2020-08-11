package file

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/stub"
	"github.com/stretchr/testify/assert"
)

var testContext = context.Background()

var fileProcessorStub = &FileProcessor{
	Config: config.Config{
		Port: "8080",
		Pubsub: config.Pubsub{
			TopicId:   "testtopic",
			ProjectId: "project",
		},
		Sample: config.Sample{
			BaseUrl: "http://localhost:8080",
		},
	},
	Client: nil,
	Ctx:    testContext,
	SampleSummary: &SampleSummary{
		Id:                            "123456",
		TotalSampleUnits:              5,
		ExpectedCollectionInstruments: 1,
	},
}

func TestScannerAndPublishSuccess(t *testing.T) {
	conn, client := stub.CreateTestPubSubServer("testtopic", testContext)
	defer conn.Close()
	defer client.Close()

	fileProcessorStub.Client = client

	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	linecount, errorCount := fileProcessorStub.Publish(scanner)

	if errorCount != 0 {
		t.Errorf("Errors have been thrown. expected: %v, actual: %v", 0, errorCount)
	}

	assert.Equal(t, 8, linecount)
}

func TestScannerAndPublishBadTopic(t *testing.T) {
	conn, client := stub.CreateTestPubSubServer("badtopic", testContext)
	defer conn.Close()
	defer client.Close()

	fileProcessorStub.Client = client

	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	linecount, errorCount := fileProcessorStub.Publish(scanner)

	if errorCount != 8 {
		t.Errorf("Invalid amount of errors thrown. expected: %v, actual: %v", 8, errorCount)
	}

	assert.Equal(t, 8, linecount)
}

func TestGetSampleSummary(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("{\"id\":\"123\",\"totalSampleUnits\":5,\"expectedCollectionInstruments\":1}"))
	}))
	ts.Start()
	defer ts.Close()
	fileProcessorStub.Config.Sample.BaseUrl = ts.URL

	sampleSummary, err := fileProcessorStub.getSampleSummary()
	assert.Nil(err, "error should be nil")
	assert.Equal("123", sampleSummary.Id, "sample summary id should match response")
	assert.Equal(5, sampleSummary.TotalSampleUnits)
	assert.Equal(1, sampleSummary.ExpectedCollectionInstruments)
}

func TestGetSampleSummaryErrors(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	ts.Start()
	defer ts.Close()
	fileProcessorStub.Config.Sample.BaseUrl = ts.URL

	sampleSummary, err := fileProcessorStub.getSampleSummary()
	assert.NotNil(err, "error should not be nil")
	assert.Nil(sampleSummary, "sample summary should be nil")
}

func TestGetCount(t *testing.T) {
	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ciCount, totalCount := fileProcessorStub.getCount(scanner)

	assert.Equal(t, 8, totalCount)
	assert.Equal(t, 1, ciCount)
}
