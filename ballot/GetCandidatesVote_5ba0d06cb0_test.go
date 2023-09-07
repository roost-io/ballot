package main

import (
	"sync"
	"testing"
)

var once sync.Once
var candidateVotesStore map[string]int

func getCandidatesVote() map[string]int {
	once.Do(func() {
		candidateVotesStore = make(map[string]int)
	})
	return candidateVotesStore
}

func TestGetCandidatesVote_5ba0d06cb0(t *testing.T) {
	candidatesVote := getCandidatesVote()
	if candidatesVote == nil {
		t.Error("Expected a non-nil map, but got nil")
	}

	// Test if the map is empty
	if len(candidatesVote) != 0 {
		t.Error("Expected an empty map, but got a map with size", len(candidatesVote))
	}

	// Test if the map works correctly by adding a value
	candidatesVote["TestCandidate"] = 5
	if len(candidatesVote) != 1 {
		t.Error("Expected a map with size 1, but got a map with size", len(candidatesVote))
	}
	if candidatesVote["TestCandidate"] != 5 {
		t.Error("Expected TestCandidate to have 5 votes, but got", candidatesVote["TestCandidate"])
	}
}
