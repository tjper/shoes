package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Error variables that occur within package.
var (
	ErrShoeIdInvalid   = errors.New("ShoeId must be set.")
	ErrShoeIdsValueDNE = errors.New("Request.Form[\"shoeIds\"] is empty.")
)

var LogErr *log.Logger

func init() {
	// Setup Logging
	file := "/log/shoes.txt"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	LogErr = log.New(f, "Error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}

// Initialize endpoint patterns and the listening port.
const (
	port = ":8080"

	trueToSizeEndpoint = "/shoes/truetosize"
)

// Init initializes and launches the shoes app.
func main() {
	// Initialize ShoesApp
	a := App{}
	a.DbClient = DbClient{Db: Db()}

	mux := http.NewServeMux()
	mux.HandleFunc(trueToSizeEndpoint, a.trueToSizeHandler)

	LogErr.Fatal(http.ListenAndServe(port, mux))

}

// App
type App struct {
	DbClient interface {
		SelectTrueToSizeByShoeId(int) ([]int, error)
		InsertTrueToSize(int, int) (int, error)
	}
}

// trueToSizeHandler routes and a truetosize request
func (a *App) trueToSizeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.getTrueToSize(w, r)
	case http.MethodPost:
		a.postTrueToSize(w, r)
	default:
		sendHttpErr(w, http.StatusNotImplemented)
	}
}

// getTrueToSize validates and processes a truetosize GET request, if successful,
// returns the true to size avg by shoe.
func (a *App) getTrueToSize(w http.ResponseWriter, r *http.Request) {
	shoeId, err := validateGetTts(r)
	if err != nil {
		sendHttpErr(w, http.StatusBadRequest)
		return
	}

	trueToSizes, err := a.DbClient.SelectTrueToSizeByShoeId(shoeId)
	if err != nil {
		sendHttpErr(w, http.StatusInternalServerError)
		return
	}

	if len(trueToSizes) == 0 {
		sendHttpErr(w, http.StatusNotFound)
		return
	}

	if err := encodeTrueToSizeAvg(w, shoeId, trueToSizes); err != nil {
		sendHttpErr(w, http.StatusInternalServerError)
		return
	}
}

// ValidateGetTts parses the GET /shoes/truetosize request, if successful,
// returns requested shoeId.
func validateGetTts(r *http.Request) (int, error) {
	var err error
	if err = r.ParseForm(); err != nil {
		LogErr.Println(err)
		return 0, err
	}

	idStr := r.FormValue("shoeId")
	if idStr == "" {
		return 0, ErrShoeIdsValueDNE
	}
	shoeId, err := strconv.Atoi(idStr)

	return shoeId, nil
}

// encodeTrueToSizeAvg encodes a ShoeId and its TrueToSize avg, if successful,
// returns nil.
func encodeTrueToSizeAvg(w http.ResponseWriter, shoeId int, trueToSizes []int) error {
	type message struct {
		ShoeId        int
		TrueToSizeAvg float64
	}
	var m message

	var sum int
	for _, v := range trueToSizes {
		sum += v
	}
	m.TrueToSizeAvg = float64(sum) / float64(len(trueToSizes))
	m.ShoeId = shoeId

	if err := json.NewEncoder(w).Encode(m); err != nil {
		LogErr.Println(err)
		return err
	}
	return nil
}

// postTrueToSize validates and processes a truetosize POST request, if successful,
// returns a http.StatusCreated.
func (a *App) postTrueToSize(w http.ResponseWriter, r *http.Request) {
	shoeId, trueToSize, err := validateShoeIdTrueToSizeJson(r.Body)
	if err != nil {
		sendHttpErr(w, http.StatusBadRequest)
		return
	}

	if _, err := a.DbClient.InsertTrueToSize(shoeId, trueToSize); err != nil {
		sendHttpErr(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// validateShoeIdTrueToSizeJson decodes and validates a ShoeId and TrueToSize set, if successful,
// returns ShoeId and TrueToSize.
func validateShoeIdTrueToSizeJson(b io.Reader) (int, int, error) {
	dec := json.NewDecoder(b)

	type e struct {
		ShoeId     int
		TrueToSize int
	}
	var m e
	for dec.More() {
		if err := dec.Decode(&m); err != nil {
			LogErr.Println(err)
			return 0, 0, err
		}
	}
	if m.ShoeId == 0 {
		return 0, 0, ErrShoeIdInvalid
	}
	if m.TrueToSize < 1 || m.TrueToSize > 5 {
		return 0, 0, ErrTrueToSizeInvalid
	}

	return m.ShoeId, m.TrueToSize, nil
}

// sendHttpErr writes an http error to the client.
func sendHttpErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
