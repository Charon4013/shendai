package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestDoubleHandler(t *testing.T) {

	testCases := []struct {
		name   string
		input  string
		result int
		status int
		err    string
	}{
		{name: "double of two", input: "4", result: 8, status: http.StatusOK, err: ""},
		{name: "double of two", input: "9", result: 18, status: http.StatusOK, err: ""},
		{name: "double of two", input: "", status: http.StatusBadRequest, err: "missing value"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/double?v="+testCase.input, nil)
			if err != nil {
				t.Fatalf("could not create a new request: %v, err: %v", request, err)
			}

			rec := httptest.NewRecorder()
			doubleHandler(rec, request)
			res := rec.Result()

			if res.StatusCode != testCase.status {
				t.Errorf("received status code %d, except %d", res.StatusCode, testCase.status)
				return
			}

			respBytes, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("cannot read all from the response body, err: %v", err)
			}
			defer res.Body.Close()

			trimedResult := strings.TrimSpace(string(respBytes))

			if res.StatusCode != http.StatusOK {
				// check the error message
				if trimedResult != testCase.err {
					t.Errorf("received status code %s, except %s", trimedResult, testCase.err)
				}
				return
			}

			// compare the returned value
			doubleVal, err := strconv.Atoi(trimedResult)
			if err != nil {
				t.Errorf("cannot convert response body to int, err: %v", err)
			}

			if doubleVal != testCase.result {
				t.Errorf("received result %d, expected %d", doubleVal, testCase.result)
			}
		})
	}

}
