package routes

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/inject"
)

func ProcessFile(w http.ResponseWriter, r *http.Request) {
	// 10MB maximum file size
	//r.ParseMultipartForm(10 << 20)
	vars := mux.Vars(r)
	sampleSummary := vars["samplesummary"]

	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.WithError(err).
			Error("Error retrieving the file")
		return
	}
	fileProcessor := inject.FileProcessor
	fileProcessor.SampleSummary = sampleSummary
	fileProcessor.ChunkCsv(file, handler)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
