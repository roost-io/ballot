package main

import (
	"testing"
)

type Vote struct {
	CandidateID string
}

var candidateVotesStore map[string]int

func saveVote(vote Vote) error {
	candidateVotesStore = getCandidatesVote()
	// add 2 votes, instead of 1
	candidateVotesStore[vote.CandidateID]++
	candidateVotesStore[vote.CandidateID]++
	return nil
}

func getCandidatesVote() map[string]int {
	if candidateVotesStore == nil {
		candidateVotesStore = make(map[string]int)
	}
	return candidateVotesStore
}

func TestSaveVote_3f5729642d(t *testing.T) {
	candidateID := "candidate1"
	vote := Vote{CandidateID: candidateID}

	err := saveVote(vote)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	if candidateVotesStore[candidateID] != 2 {
		t.Error("Expected 2 votes, got ", candidateVotesStore[candidateID])
	}

	// Test for another vote
	err = saveVote(vote)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	if candidateVotesStore[candidateID] != 4 {
		t.Error("Expected 4 votes, got ", candidateVotesStore[candidateID])
	}

	// Test for invalid candidate
	vote = Vote{CandidateID: ""}
	err = saveVote(vote)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	if candidateVotesStore[""] != 2 {
		t.Error("Expected 2 votes for invalid candidate, got ", candidateVotesStore[""])
	}
}
