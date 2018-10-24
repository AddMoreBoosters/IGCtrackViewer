package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var startTime time.Time
var tracks []igc.Track
var validID = regexp.MustCompile("^[0-9]+$")

type apiMetadata struct {
	Uptime  string
	Info    string
	Version string
}

type trackResponseData struct {
	Hdate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	TrackLength float64   `json:"track_length"`
}

func init() {
	startTime = time.Now()
}

func main() {

	urlRoot := "/igcinfo"
	myRouter := mux.NewRouter().StrictSlash(true)
	//	Disregards trailing slash, e.g. /api/ will be redirected to /api. Note that for most clients,
	//	this will turn a POST request to /api/igc/ into a GET request to /api/igc instead.
	//	Documentation: https://godoc.org/github.com/gorilla/mux#Router.StrictSlash

	myRouter.HandleFunc(urlRoot+"/api", apiInfo).Methods("GET")
	myRouter.HandleFunc(urlRoot+"/api/igc", trackRegistration).Methods("POST")
	myRouter.HandleFunc(urlRoot+"/api/igc", getAllTracks).Methods("GET")
	myRouter.HandleFunc(urlRoot+"/api/igc/{id}", getTrack).Methods("GET")
	myRouter.HandleFunc(urlRoot+"/api/igc/{id}/{field}", getTrackField).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, myRouter))

}

//	Responds with metadata about the service itself.
func apiInfo(w http.ResponseWriter, r *http.Request) {

	year, month, day, hour, min, sec := diff(startTime, time.Now())
	upTime := "P" + strconv.Itoa(year) + "Y" + strconv.Itoa(month) + "M" + strconv.Itoa(day) + "D" + "T" + strconv.Itoa(hour) + "H" + strconv.Itoa(min) + "M" + strconv.Itoa(sec) + "S"

	metadata := &apiMetadata{Uptime: upTime, Info: "Service for igc tracks.", Version: "v1"}
	json.NewEncoder(w).Encode(metadata)
}

//	Function calculates time difference
//	I copied this from https://stackoverflow.com/questions/36530251/golang-time-since-with-months-and-years
//	Reason being that the standard library's time.Since() only gives difference in hours, minutes and seconds
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

type registrationRequest struct {
	RequestURL string `json:"url"`
}

//	Gets a track from the provided url, stores it in memory, and responds with its new ID
func trackRegistration(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request registrationRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid json object", http.StatusBadRequest)
		return
	}
	reqURL, err := url.Parse(request.RequestURL)
	if err != nil {
		http.Error(w, "Invalid url", http.StatusBadRequest)
		return
	}
	track, err := igc.ParseLocation(reqURL.String())
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	tracks = append(tracks, track)
	id := len(tracks)
	json.NewEncoder(w).Encode(id)
}

//	Responds with a slice containing the ID's of all tracks stored in memory
func getAllTracks(w http.ResponseWriter, r *http.Request) {

	var ids []int

	for i := 1; i <= len(tracks); i++ {
		ids = append(ids, i)
	}
	json.NewEncoder(w).Encode(ids)

}

//	Responds with data about a single, specified track.
func getTrack(w http.ResponseWriter, r *http.Request) {

	id, err := getID(w, r)
	if err != nil {
		http.Error(w, err.Error(), id)
		return
	}

	var response = trackResponseData{
		Hdate:       tracks[id-1].Header.Date,
		Pilot:       tracks[id-1].Header.Pilot,
		Glider:      tracks[id-1].Header.GliderType,
		GliderID:    tracks[id-1].Header.GliderID,
		TrackLength: tracks[id-1].Task.Distance(),
	}

	json.NewEncoder(w).Encode(response)
}

//	Responds with the data found in a single field of a particular track.
func getTrackField(w http.ResponseWriter, r *http.Request) {

	id, err := getID(w, r)
	if err != nil {
		http.Error(w, err.Error(), id)
		return
	}

	requestedField := strings.Split(r.URL.Path, "/")[5]
	switch requestedField {
	case "pilot":
		fmt.Fprintf(w, tracks[id-1].Header.Pilot)
	case "glider":
		fmt.Fprintf(w, tracks[id-1].Header.GliderType)
	case "glider_id":
		fmt.Fprintf(w, tracks[id-1].Header.GliderID)
	case "track_length":
		fmt.Fprintf(w, strconv.FormatFloat(tracks[id-1].Task.Distance(), 'f', -1, 64))
	case "H_date":
		fmt.Fprintf(w, tracks[id-1].Header.Date.String())
	default:
		http.Error(w, "Bad request: "+requestedField+" is not a valid field.\n", 400)
		return
	}
}

//	Helper function. Gets the ID from the url and ensures it is a valid one.
func getID(w http.ResponseWriter, r *http.Request) (int, error) {
	s := strings.Split(r.URL.Path, "/")
	idString := s[4]
	if !validID.MatchString(idString) {
		return 400, errors.New("Bad request: id must be a number\n")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 500, errors.New("Internal server error: could not get id from idString\n")
	}
	if id > len(tracks) {
		return 400, errors.New("Bad request: no such id exists\n")
	}
	if id <= 0 {
		return 400, errors.New("Bad request: no such id exists\n")
	}
	return id, nil
}
