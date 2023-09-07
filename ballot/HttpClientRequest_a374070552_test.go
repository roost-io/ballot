package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func httpClientRequest(method, baseURL, path string, body io.Reader) (int, []byte, error) {
	req, err := http.NewRequest(method, baseURL+path, body)
	if err != nil {
		return 0, nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	
	return res.StatusCode, resBody, nil
}

func TestHttpClientRequest_a374070552(t *testing.T) {

	t.Run("successful request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `{"alive": true}`)
		}))
		defer ts.Close()

		status, body, err := httpClientRequest("GET", ts.URL, "/test", nil)
		if err != nil {
			t.Error("Expected no error, got ", err)
		}
		if status != 200 {
			t.Error("Expected status 200, got ", status)
		}
		if string(body) != `{"alive": true}` {
			t.Error("Expected `{'alive': true}`, got ", string(body))
		}
	})

	t.Run("failed request", func(t *testing.T) {
		status, _, err := httpClientRequest("GET", "http://invalid.url", "/test", nil)
		if err == nil {
			t.Error("Expected an error, got none")
		}
		if status != http.StatusBadRequest {
			t.Error("Expected status 400, got ", status)
		}
	})

	t.Run("request with params", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			if string(body) != `{"param": "value"}` {
				t.Error("Expected `{'param': 'value'}`, got ", string(body))
			}
			io.WriteString(w, `{"success": true}`)
		}))
		defer ts.Close()

		params := strings.NewReader(`{"param": "value"}`)
		status, body, err := httpClientRequest("POST", ts.URL, "/test", params)
		if err != nil {
			t.Error("Expected no error, got ", err)
		}
		if status != 200 {
			t.Error("Expected status 200, got ", status)
		}
		if string(body) != `{"success": true}` {
			t.Error("Expected `{'success': true}`, got ", string(body))
		}
	})
}
