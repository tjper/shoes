package shoes

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func oneShoeJson() io.Reader {
	return strings.NewReader(`{"shoeId":1}`)
}

func twoShoesJson() io.Reader {
	return strings.NewReader(`{"shoeId":1}{"shoeId":2}`)
}

func Test_getTrueToSize(t *testing.T) {
	type test struct {
		name           string
		expectedStatus int
		expectedBody   string
		req            *http.Request
	}

	tests := []test{
		test{"No body.", http.StatusBadRequest, "Bad Request\n", httptest.NewRequest(http.MethodGet, trueToSizeEndpoint, nil)},
		test{"One shoe.", http.StatusOK, "{\"ShoeId\":1,\"TrueToSize\":2}\n", httptest.NewRequest(http.MethodGet, trueToSizeEndpoint, oneShoeJson())},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			a := App{}
			a.DbClient = mockDbClient{}
			w := httptest.NewRecorder()

			a.getTrueToSize(w, te.req)
			resp := w.Result()
			if resp.StatusCode != te.expectedStatus {
				t.Errorf("Expected Status Code = %v, Actual Status Code = %v\n", te.expectedStatus, resp.StatusCode)
			}

			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) != te.expectedBody {
				t.Errorf("Expected body = %v, Actual body = %v,\n", te.expectedBody, string(body))
			}
		})
	}
}

func Test_decodeShoesIdsJson(t *testing.T) {
	type test struct {
		name        string
		body        io.Reader
		expectedRes []int
		expectedErr error
	}

	tests := []test{
		test{"No body.", strings.NewReader(""), nil, ErrEmptyReqBody},
		test{"One Shoe.", oneShoeJson(), []int{1}, nil},
		test{"Two Shoes.", twoShoesJson(), []int{1, 2}, nil},
	}
	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res, err := decodeShoeIdsJson(te.body)
			for i, v := range res {
				if te.expectedRes[i] != v {
					t.Errorf("Expected Res = %v, Actual Res = %v\n", te.expectedRes, res)
					break
				}
			}
			if err != te.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", te.expectedErr, err)
			}
		})
	}
}

/*
func test_encodeShoesTrueToSizes(t *testing.T) {
	type test struct {
		name             string
		shoesTrueToSizes [int][]int
		expectedErr      error
	}
}
*/

func Test_postTrueToSize(t *testing.T) {
	type test struct {
		name           string
		expectedStatus int
		req            *http.Request
	}

	tests := []test{
		test{"No body.", http.StatusBadRequest, httptest.NewRequest(http.MethodGet, trueToSizeEndpoint, nil)},
		test{"One shoe.", http.StatusOK, httptest.NewRequest(http.MethodGet, trueToSizeEndpoint, strings.NewReader(`{"ShoeId": 1, "TrueToSize":[1,2,3]}`))},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			a := App{}
			a.DbClient = mockDbClient{}
			w := httptest.NewRecorder()

			a.getTrueToSize(w, te.req)
			resp := w.Result()
			if resp.StatusCode != te.expectedStatus {
				t.Errorf("Expected Status Code = %v, Actual Status Code = %v\n", te.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// mockDbClient for testing.
type mockDbClient struct{}

func (m mockDbClient) SelectTrueToSizeByShoesId(shoeIds []int) (map[int][]int, error) {
	switch len(shoeIds) {
	case 1:
		return shoeSizes, nil
	case 2:
		return shoesSizes, nil
	}
	return nil, ErrZeroShoesIds
}
func (m mockDbClient) InsertTrueToSizes(ShoesTrueToSize map[int][]int) (int, error) { return 0, nil }
