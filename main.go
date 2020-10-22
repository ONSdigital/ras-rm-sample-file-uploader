package main

import (
	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"net/http"

	"github.com/gorilla/mux"
)

var logger *zap.Logger

func init() {
	logger, _ = zapdriver.NewProduction()
	defer logger.Sync()
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/fileupload", routes.ProcessFile)
	router.HandleFunc("/info", routes.Info)
	http.ListenAndServe(":8080", router)
}
