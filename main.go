package main

import (
	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/fileupload", routes.ProcessFile)
	router.HandleFunc("/info", routes.Info)
	http.ListenAndServe(":8080", router)
}
