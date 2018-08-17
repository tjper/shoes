// Package shoes opens an API to interact with the shoes resource.
package shoes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Error variables that occur within package.
var (
	ErrTrueSizeNotInRange = errors.New("TrueToSize must be an int between 1 and 5 inclusive.")
	ErrEmptyReqBody       = errors.New("Request body is empty.")
)

// App
type App struct {
	DbClient interface {
		SelectTrueToSizeByShoesId([]int) (map[int][]int, error)
		InsertTrueToSizes(map[int][]int) (int, error)
	}
}

// Initialize endpoint patterns and the listening port.
const (
	port = ":8080"

	trueToSizeEndpoint = "/shoes/truetosize"
)

// Init initializes and launches the shoes app.
func (a *App) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc(trueToSizeEndpoint, a.trueToSizeHandler)

	LogErr.Fatal(http.ListenAndServe(port, mux))
}

// trueToSizeHandler routes and a truetosize request
func (a *App) trueToSizeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.getTrueToSize(w, r)
	case http.MethodPost:
		a.postTrueToSize(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

// getTrueToSize validates and processes a truetosize GET request, if successful,
// returns the true to size avg by shoe.
func (a *App) getTrueToSize(w http.ResponseWriter, r *http.Request) {
	shoeIds, err := decodeShoeIdsJson(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	shoesTrueToSizes, err := a.DbClient.SelectTrueToSizeByShoesId(shoeIds)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := encodeShoesTrueToSizes(w, shoesTrueToSizes); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
}

// decodeShoeIdsJson decodes a ShoeId Json set, if successful,
// returns a slice of ShoeIds.
func decodeShoeIdsJson(b io.Reader) ([]int, error) {
	dec := json.NewDecoder(b)

	type shoe struct {
		ShoeId int
	}
	shoeIds := make([]int, 0)

	for dec.More() {
		var s shoe
		if err := dec.Decode(&s); err != nil {
			LogErr.Println(err)
			return nil, err
		}
		shoeIds = append(shoeIds, s.ShoeId)
	}
	if len(shoeIds) == 0 {
		LogErr.Println(ErrEmptyReqBody)
		return nil, ErrEmptyReqBody
	}
	return shoeIds, nil
}

// encodeShoesTrueToSizes encodes a ShoeId to TrueToSize sets, if successful,
// returns nil.
func encodeShoesTrueToSizes(w http.ResponseWriter, shoesTrueToSizes map[int][]int) error {
	type message struct {
		ShoeId     int
		TrueToSize float64
	}
	mSet := make([]message, 0)

	for id, tts := range shoesTrueToSizes {
		var sum int
		for _, v := range tts {
			sum += v
		}
		avg := float64(sum) / float64(len(tts))
		mSet = append(mSet, message{id, avg})
	}

	enc := json.NewEncoder(w)
	for _, m := range mSet {
		if err := enc.Encode(m); err != nil {
			LogErr.Println(err)
			return err
		}
	}
	return nil
}

// postTrueToSize validates and processes a truetosize POST request, if successful,
// returns a http.StatusCreated.
func (a *App) postTrueToSize(w http.ResponseWriter, r *http.Request) {
	shoesTrueToSizes, err := decodeShoesTrueToSizesJson(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if _, err := a.DbClient.InsertTrueToSizes(shoesTrueToSizes); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// decodeShoesTrueToSizesJson decodes ShoeId to TrueToSize sets, if successful,
// returns a mapping of ShoeId to TrueToSize sets.
func decodeShoesTrueToSizesJson(b io.Reader) (map[int][]int, error) {
	dec := json.NewDecoder(b)
	if !dec.More() {
		LogErr.Println(ErrEmptyReqBody)
		return nil, ErrEmptyReqBody
	}

	type message struct {
		ShoeId     int
		TrueToSize []int
	}

	shoesTrueToSizes := make(map[int][]int)
	for dec.More() {
		var m message
		if err := dec.Decode(&m); err != nil {
			LogErr.Println(err)
			return nil, err
		}

		for _, s := range m.TrueToSize {
			if s < 1 || s > 5 {
				LogErr.Println(ErrTrueSizeNotInRange)
				return nil, ErrTrueSizeNotInRange
			}
		}

		if _, exists := shoesTrueToSizes[m.ShoeId]; exists {
			shoesTrueToSizes[m.ShoeId] = append(shoesTrueToSizes[m.ShoeId], m.TrueToSize...)

		} else {
			shoesTrueToSizes[m.ShoeId] = m.TrueToSize
		}
	}
	return shoesTrueToSizes, nil
}
