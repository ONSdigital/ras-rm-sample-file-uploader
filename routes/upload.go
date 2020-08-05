package routes

import (
	"fmt"
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
		//w.WriteHeader(http.StatusBadRequest)
		//return
	}
	err = inject.FileProcessor.ChunkCsv(file, handler)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
