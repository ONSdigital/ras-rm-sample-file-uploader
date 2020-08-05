package main

import (
	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/fileupload", routes.ProcessFile)
	http.ListenAndServe(":8080", router)
}
