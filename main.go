package main

import (
	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	"github.com/spf13/viper"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	viper.AutomaticEnv()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/fileupload", routes.ProcessFile)
	http.ListenAndServe(":8080", router)
}
