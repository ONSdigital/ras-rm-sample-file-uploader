package main

import (
	"net/http"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	//TODO this can be reintroduced once we remove the sample summary creation from sample service
	//router.HandleFunc("/samples/{type}/fileupload", routes.ProcessFile)
	router.HandleFunc("/sample-summary/{samplesummary}/samples/fileupload", routes.ProcessFile)
	http.ListenAndServe(":8080", router)
}
