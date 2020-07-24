package file

import (
	"bufio"
	"context"
	"os"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/stub"
	"testing"
	"github.com/stretchr/testify/assert"
)

var testContext = context.Background()

var fileProcessorStub = &FileProcessor{
	Config: config.Config{
		Port: "8080",
		Pubsub: config.Pubsub{
			TopicId: "testtopic",
			ProjectId: "project",
	    },
	},
	Client: nil,
	Ctx: testContext,
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

	errorCount := fileProcessorStub.Publish(scanner)

	if errorCount != 0 {
		t.Errorf("Errors have been thrown. expected: %v, actual: %v", 0, errorCount)
	}
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

	errorCount := fileProcessorStub.Publish(scanner)

	if errorCount != 8 {
		t.Errorf("Invalid amount of errors thrown. expected: %v, actual: %v", 8, errorCount)
	}
}