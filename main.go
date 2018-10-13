package main 

import (
	"fmt"
	"time"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

var startTime time.Time
var tracks []igc.Track

type ApiMetadata struct {
	Uptime	string
	Info 	string 
	Version	string
}

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

func apiInfo(w http.ResponseWriter, r *http.Request) {

	//	The time package's Duration type only shows hours, minutes and seconds, so to get the full
	//	format I use Date instead. The side effect is that since 0000.00.00 is not a valid date,
	//	uptimes of less than 1 day become -0001.11.30 instead. Oh well.
	exampleTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Add(time.Since(startTime))

	metadata := &ApiMetadata{Uptime: exampleTime.String(), Info: "Service for igc tracks.", Version: "v1"}
	json.NewEncoder(w).Encode(metadata)
}

func trackRegistration(w http.ResponseWriter, r *http.Request) {

	url, err := ioutil.ReadAll(r.Body)
	if (err != nil) {
		fmt.Errorf("Problem reading the url", err)
	}
	track, err := igc.ParseLocation(string(url))
	if (err != nil) {
		fmt.Errorf("Problem reading the track", err)
	}
	tracks = append(tracks, track)
	id := len(tracks)
	json.NewEncoder(w).Encode(id)
}

func getAllTracks(w http.ResponseWriter, r *http.Request) {

	var ids []int

	for i := 1; i <= len(tracks); i++ {
		ids = append(ids, i)
	}
	json.NewEncoder(w).Encode(ids)

}

func getTrack(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Endpoint hit: getTrack\n")
	http.Error(w, "Not implemented\n", 501)
}

func getTrackField(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Endpoint hit: getTrackField\n")
	http.Error(w, "Not implemented\n", 501)
}
