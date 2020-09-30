package routes

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/inject"
)

func ProcessFile(w http.ResponseWriter, r *http.Request) {
	// 10MB maximum file size
	//r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.WithError(err).
			Error("Error retrieving the file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sampleSummary, err := inject.FileProcessor.ChunkCsv(file, handler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	js, err := json.Marshal(sampleSummary)
	log.WithField("json", string(js)).Info("returning sample summary")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//Info endpoint handler returns info like name, version, origin, commit, branch
//and built
func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
