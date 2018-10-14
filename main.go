package main 

import (
	"fmt"
	"time"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"strings"
	"os"
	"strconv"
	"errors"
	"regexp"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

var startTime time.Time
var tracks []igc.Track
var validID = regexp.MustCompile("^[0-9]+$")

type ApiMetadata struct {
	Uptime	string
	Info 	string 
	Version	string
}

type TrackResponseData struct {
	Hdate time.Time 	`json:"H_date"`
	Pilot string 		`json:"pilot"`
	Glider string 		`json:"glider"`
	GliderID string 	`json:"glider_id"`
	TrackLength float64 `json:"track_length"`
}

func init () {
	startTime = time.Now()
}

func main () {
	
	urlRoot := "/igcinfo"
	myRouter := mux.NewRouter().StrictSlash(true)	
	//	Disregards trailing slash, e.g. /api/ will be redirected to /api. Note that for
	//	most clients, this will turn a POST request to /api/igc/ into a GET request to /api/igc
	//	Documentation: https://godoc.org/github.com/gorilla/mux#Router.StrictSlash

	myRouter.HandleFunc(urlRoot + "/api", apiInfo).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc", trackRegistration).Methods("POST")
	myRouter.HandleFunc(urlRoot + "/api/igc", getAllTracks).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc/{id}", getTrack).Methods("GET")
	myRouter.HandleFunc(urlRoot + "/api/igc/{id}/{field}", getTrackField).Methods("GET")

	port := os.Getenv("PORT")
	if (port == "") {
		log.Fatal("Port must be set")
	}

	log.Fatal(http.ListenAndServe(":" + port, myRouter))
	
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

	id, err := getID(w, r)
	if (err != nil) {
		http.Error(w, err.Error(), id)
		return
	}

	var response = TrackResponseData {
		Hdate: tracks[id - 1].Header.Date,
		Pilot: tracks[id - 1].Header.Pilot,
		Glider: tracks[id - 1].Header.GliderType,
		GliderID: tracks[id - 1].Header.GliderID,
		TrackLength: tracks[id - 1].Task.Distance(),
	}

	json.NewEncoder(w).Encode(response)
}

func getTrackField(w http.ResponseWriter, r *http.Request) {
	
	id, err := getID(w, r)
	if (err != nil) {
		http.Error(w, err.Error(), id)
		return
	}

	requestedField := strings.Split(r.URL.Path, "/")[5]
	switch requestedField {
	case "pilot":
		fmt.Fprintf(w, tracks[id - 1].Header.Pilot)
	case "glider":
		fmt.Fprintf(w, tracks[id - 1].Header.GliderType)
	case "glider_id":
		fmt.Fprintf(w, tracks[id - 1].Header.GliderID)
	case "track_length":
		fmt.Fprintf(w, strconv.FormatFloat(tracks[id - 1].Task.Distance(), 'f', -1, 64))
	case "H_date":
		fmt.Fprintf(w, tracks[id - 1].Header.Date.String())
	default:
		http.Error(w, "Bad request: " + requestedField + " is not a valid field.\n", 400)
		return
	}
}

func getID (w http.ResponseWriter, r *http.Request) (int, error) {
	s := strings.Split(r.URL.Path, "/")
	idString := s[4]
	if (!validID.MatchString(idString)) {
		return 400, errors.New("Bad request: id must be a number\n")
	}
	id, err := strconv.Atoi(idString)
	if (err != nil) {
		return 500, errors.New("Internal server error: could not get id from idString\n")
	}
	if (id > len(tracks)) {
		return 400, errors.New("Bad request: no such id exists\n")
	}
	return id, nil
}
