package main 

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func main () {

	urlRoot := "/igcinfo"
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc(urlRoot + "/api", apiInfo).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc", trackRegistration).Methods("POST")
	myRouter.HandleFunc(urlRoot + "/api/igc", getAllTracks).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc/<id>", getTrack).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc/<id>/<field>", getTrackField).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func apiInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint hit: apiInfo")
}

func trackRegistration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint hit: trackRegistration")
}

func getAllTracks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint hit: getAllTracks")
}

func getTrack(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint hit: getTrack")
}

func getTrackField(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint hit: getTrackField")
}