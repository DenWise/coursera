package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	ID string
	Response string
	StatusCode int
}

func TestGetUser(t *testing.T) {
	cases := []TestCase{
		TestCase{
			ID: "42",
			Response: `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			ID: "500",
			Response: `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusInternalServerError,
		},
	}

	for caseNum, item := range cases {
		url := "http://gagaga.com/api/user?id=" + item.ID
		req := httptest.NewRequest(http.MethodGet, url, nil) // mock request
		w := httptest.NewRecorder() // response writer for GetUser

		GetUser(w, req) // mock execute GetUser handler

		if w.Code != item.StatusCode {
			t.Errorf("[%d] Incorrect status code: got %d, expected %d", caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result() // getting response from mock response writer
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] Incorrect response body: got %s, expected %s", caseNum, bodyStr, item.Response)
		}
	}
}