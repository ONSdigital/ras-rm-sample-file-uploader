package routes

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/file"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/inject"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/stub"
	"github.com/stretchr/testify/assert"
)

var fileProcessorStub file.FileProcessor
var ctx = context.Background()

func init() {
	_, client := stub.CreateTestPubSubServer("testtopic", ctx)
	fileProcessorStub = file.FileProcessor{
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
		Client: client,
		Ctx:    ctx,
	}
}

func TestFileUploadSuccess(t *testing.T) {
	expectedReturn := "{\"id\":\"123\",\"totalSampleUnits\":5,\"expectedCollectionInstruments\":1}"
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(expectedReturn))
	}))
	ts.Start()
	defer ts.Close()
	fileProcessorStub.Config.Sample.BaseUrl = ts.URL

	inject.FileProcessor = fileProcessorStub
	path := "../file/sample_test_file.csv"
	file, err := os.Open(path)
	assert.Nil(t, err)

	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	assert.Nil(t, err)
	io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest("POST", "/samples/B/fileupload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	ProcessFile(res, req)

	assert.Equal(t, 202, res.Code)

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, expectedReturn, bodyString)
}
