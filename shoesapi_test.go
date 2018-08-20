package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const (
	host = "http://localhost:8080"
)

func Test_trueToSizeHandler(t *testing.T) {
	tests := []struct {
		name           string
		req            *http.Request
		expectedStatus int
	}{
		{"GET request", httptest.NewRequest(http.MethodGet, host+"/shoes/truetosize?shoeId=1", nil), http.StatusOK},
		{"POST request", httptest.NewRequest(http.MethodPost, host+"/shoes/truetosize", bytes.NewReader([]byte("{\"ShoeId\":1, \"TrueToSize\":2}"))), http.StatusCreated},
		{"PUT request", httptest.NewRequest(http.MethodPut, host+"/shoes/truetosize", nil), http.StatusNotImplemented},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := App{DbClient: mockDbClient{}}
			w := httptest.NewRecorder()

			a.trueToSizeHandler(w, test.req)
			resp := w.Result()
			if resp.StatusCode != test.expectedStatus {
				notExpected(t, test.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func Test_validateGetTts(t *testing.T) {
	tests := []struct {
		name        string
		req         *http.Request
		expectedRes int
		expectedErr error
	}{
		{"No values.", httptest.NewRequest(http.MethodGet, host+"/shoes/truetosize", nil), 0, ErrShoeIdsValueDNE},
		{"One shoeId.", httptest.NewRequest(http.MethodGet, host+"/shoes/truetosize?shoeId=1", nil), 1, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := validateGetTts(test.req)
			if err != test.expectedErr {
				notExpected(t, test.expectedErr, err)
			}
			if !reflect.DeepEqual(res, test.expectedRes) {
				notExpected(t, test.expectedRes, res)
			}

		})
	}
}

func Test_encodeTrueToSizeAvg(t *testing.T) {
	type e struct {
		ShoeId        int
		TrueToSizeAvg float64
	}
	tests := []struct {
		name         string
		shoeId       int
		trueToSizes  []int
		expectedBody e
		expectedErr  error
	}{
		{"One shoe, one truetosize.", 1, []int{2}, e{1, 2}, nil},
		{"One shoe, multiple truetosize.", 1, []int{1, 2}, e{1, 1.5}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := encodeTrueToSizeAvg(w, test.shoeId, test.trueToSizes)
			if err != test.expectedErr {
				notExpected(t, test.expectedErr, err)
			}
			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			expectedBody, err := json.Marshal(test.expectedBody)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Equal(body, expectedBody) {
				notExpected(t, expectedBody, body)
			}
		})
	}
}

func Test_validateShoeIdTrueToSizeJson(t *testing.T) {
	tests := []struct {
		name               string
		json               string
		expectedShoeId     int
		expectedTrueToSize int
		expectedErr        error
	}{
		{"No body.", "", 0, 0, ErrShoeIdInvalid},
		{"No shoe, one trueToSize.", `{"truetosize":3}`, 0, 0, ErrShoeIdInvalid},
		{"One Shoe, no trueToSize.", `{"ShoeId":1}`, 0, 0, ErrTrueToSizeInvalid},
		{"One shoe, one truetosize.", `{"shoeId": 1, "truetosize": 2}`, 1, 2, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shoeId, trueToSize, err := validateShoeIdTrueToSizeJson(strings.NewReader(test.json))
			if err != test.expectedErr {
				notExpected(t, test.expectedErr, err)
			}
			if shoeId != test.expectedShoeId {
				notExpected(t, test.expectedShoeId, shoeId)
			}
			if trueToSize != test.expectedTrueToSize {
				notExpected(t, test.expectedTrueToSize, trueToSize)
			}
		})
	}
}

func notExpected(t *testing.T, expected, actual interface{}) {
	t.Errorf("Expected %T = %v, Actual %T = %v\n", expected, expected, actual, actual)
}

// mockDbClient for testing.
type mockDbClient struct{}

func (m mockDbClient) SelectTrueToSizeByShoeId(shoeId int) ([]int, error)   { return []int{1}, nil }
func (m mockDbClient) InsertTrueToSize(shoeId, trueToSize int) (int, error) { return 1, nil }
