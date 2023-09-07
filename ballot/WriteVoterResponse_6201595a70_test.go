package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Status struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func writeVoterResponse(w http.ResponseWriter, status Status) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(status)
	if err != nil {
		log.Println("error marshaling response to vote request. error: ", err)
		return
	}
	w.Write(resp)
}

func TestWriteVoterResponse_6201595a70(t *testing.T) {
	t.Run("Expected Response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		status := Status{Message: "Vote Successful", Code: 200}
		writeVoterResponse(recorder, status)
		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)
		result.Body.Close()
		if result.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", result.StatusCode)
		}
		var resp Status
		json.Unmarshal(body, &resp)
		if resp.Message != status.Message {
			t.Errorf("Expected message %q; got %q", status.Message, resp.Message)
		}
		if resp.Code != status.Code {
			t.Errorf("Expected code %v; got %v", status.Code, resp.Code)
		}
	})

	t.Run("Marshaling Error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		status := Status{Message: string([]byte{0x80, 0x81}), Code: 200}
		writeVoterResponse(recorder, status)
		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)
		result.Body.Close()
		if result.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", result.StatusCode)
		}
		if !bytes.Equal(body, []byte{}) {
			t.Errorf("Expected empty body; got %q", body)
		}
	})
}
