package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
	"fmt"
)

type Status struct {
	Message string
	Code    int
}

func runTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	log.Println("ballot endpoint tests running")
	status := Status{}
	err := testBallot()
	if err != nil {
		status.Message = fmt.Sprintf("Test Cases Failed with error : %v", err)
		status.Code = http.StatusBadRequest
	} else {
		status.Message = "Test Cases passed"
		status.Code = http.StatusOK
	}
	writeVoterResponse(w, status)
}

func testBallot() error {
	// TODO: Implement the test logic here
	return nil
}

func writeVoterResponse(w http.ResponseWriter, status Status) {
	// TODO: Implement the response writing logic here
}

func TestRunTest_a2bb790205(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(runTest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"Message":"Test Cases passed","Code":200}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRunTest_a2bb790205_Failure(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(runTest)

	// TODO: Make TestBallot return an error to simulate failure
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"Message":"Test Cases Failed with error : ","Code":400}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
