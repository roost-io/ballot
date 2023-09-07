package main

import (
	"testing"
	"sort"
)

type ResultBoard struct {
	Results    []CandidateVotes
	TotalVotes int
}

type CandidateVotes struct {
	CandidateID string
	Votes       int
}

func getCandidatesVote() map[string]int {
	return map[string]int{
		"candidate1": 5,
		"candidate2": 3,
		"candidate3": 7,
	}
}

func countVote() (res ResultBoard, err error) {
	votes := getCandidatesVote()
	for candidateID, vote := range votes {
		res.Results = append(res.Results, CandidateVotes{candidateID, vote})
		res.TotalVotes += vote
	}

	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].Votes > res.Results[j].Votes
	})
	return res, err
}

func TestCountVote_e644ee2864(t *testing.T) {
	res, err := countVote()
	if err != nil {
		t.Error("Error occurred: ", err)
	}

	if res.TotalVotes != 15 {
		t.Error("Total votes count mismatch")
	}

	if len(res.Results) != 3 {
		t.Error("Candidates count mismatch")
	}

	if res.Results[0].CandidateID != "candidate3" || res.Results[0].Votes != 7 {
		t.Error("Votes count for candidate3 mismatch")
	}

	if res.Results[1].CandidateID != "candidate1" || res.Results[1].Votes != 5 {
		t.Error("Votes count for candidate1 mismatch")
	}

	if res.Results[2].CandidateID != "candidate2" || res.Results[2].Votes != 3 {
		t.Error("Votes count for candidate2 mismatch")
	}
}
