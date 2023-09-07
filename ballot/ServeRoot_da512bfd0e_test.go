package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeRoot_da512bfd0e(t *testing.T) {
	// Test case 1: GET request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serveRoot)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Test case 2: POST request with valid data
	vote := `{"VoterID": "123", "CandidateID": "456"}`
	req, err = http.NewRequest("POST", "/", strings.NewReader(vote))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Test case 3: POST request with invalid data
	vote = `{"VoterID": "", "CandidateID": "456"}`
	req, err = http.NewRequest("POST", "/", strings.NewReader(vote))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Test case 4: Unsupported request method
	req, err = http.NewRequest("PUT", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
