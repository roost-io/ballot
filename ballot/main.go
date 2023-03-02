package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strings"
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

// ResultBoard to send when voting result requested
type ResultBoard struct {
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
	// add 2 votes, instead of 1
	candidateVotesStore[vote.CandidateID]++
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

func TestBallot() error {
	_, result, err := httpClientRequest(http.MethodGet, net.JoinHostPort("", port), "/", nil)
	if err != nil {
		log.Printf("Failed to get ballot count resp:%s error:%+v", string(result), err)
		return err
	}
	log.Println("get ballot resp:", string(result))
	var initalRespData ResultBoard
	if err = json.Unmarshal(result, &initalRespData); err != nil {
		log.Printf("Failed to unmarshal get ballot response. %+v", err)
		return err
	}

	var ballotvotereq Vote
	ballotvotereq.CandidateID = fmt.Sprint(rand.Intn(10))
	ballotvotereq.VoterID = fmt.Sprint(rand.Intn(10))
	reqBuff, err := json.Marshal(ballotvotereq)
	if err != nil {
		log.Printf("Failed to marshall post ballot request %+v", err)
		return err
	}

	_, result, err = httpClientRequest(http.MethodPost, net.JoinHostPort("", port), "/", bytes.NewReader(reqBuff))
	if err != nil {
		log.Printf("Failed to get ballot count resp:%s error:%+v", string(result), err)
		return err
	}
	log.Println("post ballot resp:", string(result))
	var postballotResp Status
	if err = json.Unmarshal(result, &postballotResp); err != nil {
		log.Printf("Failed to unmarshal post ballot response. %+v", err)
		return err
	}
	if postballotResp.Code != 201 {
		return errors.New("post ballot resp status code")
	}

	_, result, err = httpClientRequest(http.MethodGet, net.JoinHostPort("", port), "/", nil)
	if err != nil {
		log.Printf("Failed to get final ballot count resp:%s error:%+v", string(result), err)
		return err
	}
	log.Println("get final ballot resp:", string(result))
	var finalRespData ResultBoard
	if err = json.Unmarshal(result, &finalRespData); err != nil {
		log.Printf("Failed to unmarshal get final ballot response. %+v", err)
		return err
	}
	if finalRespData.TotalVotes-initalRespData.TotalVotes != 1 {
		return errors.New("ballot vote count error")
	}
	return nil
}

func httpClientRequest(operation, hostAddr, command string, params io.Reader) (int, []byte, error) {

	url := "http://" + hostAddr + command
	if strings.Contains(hostAddr, "http://") {
		url = hostAddr + command
	}

	req, err := http.NewRequest(operation, url, params)
	if err != nil {
		return http.StatusBadRequest, nil, errors.New("Failed to create HTTP request." + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	defer resp.Body.Close()

	body, ioErr := ioutil.ReadAll(resp.Body)
	if hBit := resp.StatusCode / 100; hBit != 2 && hBit != 3 {
		if ioErr != nil {
			ioErr = fmt.Errorf("status code error %d", resp.StatusCode)
		}
	}
	return resp.StatusCode, body, ioErr
}

// countVote to show in result board.
func countVote() (res ResultBoard, err error) {
	votes := getCandidatesVote()
	for candidateID, votes := range votes {
		res.Results = append(res.Results, CandidateVotes{candidateID, votes})
		res.TotalVotes += votes
	}

	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].Votes > res.Results[j].Votes
	})
	return res, err
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodGet:
		defer r.Body.Close()
		log.Println("result request received")

		voteData, err := countVote()
		out, err := json.Marshal(voteData)
		if err != nil {
			log.Println("error marshaling response to result request. error: ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(out)

	case http.MethodPost:
		log.Println("vote received")
		vote := Vote{}
		status := Status{}
		defer writeVoterResponse(w, status)
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&vote)
		if err != nil || vote.CandidateID == "" {
			log.Printf("error parsing vote data. error: %v\n", err)
			status.Code = http.StatusBadRequest
			status.Message = "Vote is not valid. Vote can not be saved"
			return
		}
		log.Printf("Voting done by voter: %s to candidate: %s\n", vote.VoterID, vote.CandidateID)
		err = saveVote(vote)
		if err != nil {
			log.Println(err)
			status.Code = http.StatusBadRequest
			status.Message = "Vote is not valid. Vote can not be saved"
			return
		}
		status.Code = http.StatusCreated
		status.Message = "Vote saved sucessfully"
		return

	default:
		status := Status{}
		status.Code = http.StatusMethodNotAllowed
		status.Message = "Bad Request. Vote can not be saved"
		return
	}
}

// runTest add votes and also checks the validity
func runTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	log.Println("ballot endpoint tests running")
	status := Status{}
	err := TestBallot()
	if err != nil {
		status.Message = fmt.Sprintf("Test Cases Failed with error : %v", err)
		status.Code = http.StatusBadRequest
	}
	status.Message = "Test Cases passed"
	status.Code = http.StatusOK
	writeVoterResponse(w, status)
}

func main() {
	log.Println("ballot is ready to store votes")
	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/tests/run", runTest)
	log.Println(http.ListenAndServe(net.JoinHostPort("", port), nil))
}
