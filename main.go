package main

import (
	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	viper.AutomaticEnv()
	configureLogging()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/fileupload", routes.ProcessFile)
	http.ListenAndServe(":8080", router)
}

func configureLogging() {
	verbose := viper.GetBool("VERBOSE")
	log.SetFormatter(&log.JSONFormatter{})
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}