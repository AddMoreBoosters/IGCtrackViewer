package main 

import (
	"fmt"
	"time"
	"net/http"
	"log"
	"encoding/json"
	//"ioutil"
	"github.com/gorilla/mux"
)

//	TODO create internal storage system. Separate package?

var startTime time.Time

func init () {
	startTime = time.Now()
}

func main () {
	
	urlRoot := "/igcinfo"
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc(urlRoot + "/api", apiInfo).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc", trackRegistration).Methods("POST")
	myRouter.HandleFunc(urlRoot + "/api/igc", getAllTracks).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc/{id}", getTrack).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc/{id}/{field}", getTrackField).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
	
}

type ApiMetadata struct {
	Uptime	string
	Info 	string 
	Version	string
}

func apiInfo(w http.ResponseWriter, r *http.Request) {

	//	The time package's Duration type only shows hours, minutes and seconds, so to get the full
	//	format I use Date instead. The side effect is that since 0000.00.00 is not a valid date,
	//	uptimes of less than 1 day become -0001.11.30 instead. Oh well.
	exampleTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Add(time.Since(startTime))

	metadata := &ApiMetadata{Uptime: exampleTime.String(), Info: "Service for igc tracks.", Version: "v1"}
	json.NewEncoder(w).Encode(metadata)
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
