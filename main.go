package main

import (
	"net/http"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/{type}/fileupload", routes.ProcessFile)
	http.ListenAndServe(":8080", router)
}
