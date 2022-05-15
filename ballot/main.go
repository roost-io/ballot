package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sort"
	"sync"
)

var port string = "8080"
var once sync.Once

// candidateVotesStore holds map[candidate_id] = vote_count
var candidateVotesStore map[string]int

// Vote data
type Vote struct {
	CandidateID string `json:"candidate_id"`
	VoterID     string `json:"voter_id"`
}

// CandidateVotes contains candidates and their vote counts
type CandidateVotes struct {
	CandidateID string `json:"candidate_id"`
	Votes       int    `json:"vote_count"`
}

// Response to send when voting result requested
type Response struct {
	Results    []CandidateVotes `json:"results"`
	TotalVotes int              `json:"total_votes"`
}

// Status to be sent in response to API request
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// getVote returns empty data instead of nil if voting not happened
func getCandidatesVote() map[string]int {
	once.Do(func() {
		candidateVotesStore = make(map[string]int)
	})
	return candidateVotesStore
}

// saveVote regardless of who voted
func saveVote(vote Vote) error {
	candidateVotesStore = getCandidatesVote()
	candidateVotesStore[vote.CandidateID]++
	return nil
}

func writeVoterResponse(w http.ResponseWriter, status Status) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(status)
	if err != nil {
		log.Println("error marshaling response to vote request. error: ", err)
	}
	w.Write(resp)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodGet:
		defer r.Body.Close()
		log.Println("result request received")
		res := Response{}
		votes := getCandidatesVote()
		for candidateID, votes := range votes {
			res.Results = append(res.Results, CandidateVotes{candidateID, votes})
			res.TotalVotes += votes
		}

		sort.Slice(res.Results, func(i, j int) bool {
			return res.Results[i].Votes > res.Results[j].Votes
		})

		log.Printf("result data: %+v", res)

		out, err := json.Marshal(res)
		if err != nil {
			log.Println("error marshaling response to result request. error: ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)

	case http.MethodPost:
		log.Println("vote received")
		vote := Vote{}
		status := Status{}
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&vote)
		if err != nil {
			log.Printf("error parsing vote data. error: %v\n", err)
			status.Code = http.StatusBadRequest
			status.Message = "Vote is not valid. Vote can not be saved"
			writeVoterResponse(w, status)
			return
		}
		log.Printf("Voting done by voter: %s to candidate: %s\n", vote.VoterID, vote.CandidateID)
		err = saveVote(vote)
		if err != nil {
			log.Println(err)
			status.Code = http.StatusBadRequest
			status.Message = "Vote is not valid. Vote can not be saved"
			writeVoterResponse(w, status)
			return
		}
		status.Code = http.StatusCreated
		status.Message = "Vote saved sucessfully"
		writeVoterResponse(w, status)
		return

	default:
		status := Status{}
		status.Code = http.StatusMethodNotAllowed
		status.Message = "Bad Request. Vote can not be saved"
		writeVoterResponse(w, status)
		return
	}
}

func main() {
	log.Println("ballot is ready to store votes")
	http.HandleFunc("/", serveRoot)
	log.Println(http.ListenAndServe(net.JoinHostPort("", port), nil))
}
