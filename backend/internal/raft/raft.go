package raft

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type NodeRole int

const (
	Follower NodeRole = iota
	Candidate
	Leader
)

type Log struct {
	command string
	arg1    string
	arg2    string
}

type Raft struct {
	members         []string
	follower        []string
	leader          string
	self            string
	log             []Log
	mu              sync.Mutex
	role            NodeRole
	term            int
	votedFor        string
	votes           int
	electionTimeout *time.Timer
}

type SubscribeReq struct {
	IPAddress string `json:"IPAddress"`
}

type Response struct {
	Status     string   `json:"status"`
	Message    string   `json:"message"`
	RaftMember []string `json:"raftMember"`
}

var raft = Raft{
	members:  []string{},
	leader:   "",
	log:      []Log{},
	role:     Follower,
	term:     0,
	votedFor: "",
}

func AddLog(cmd string, arg1 string, arg2 string) {
	log := Log{
		command: cmd,
		arg1:    arg1,
		arg2:    arg2,
	}
	raft.log = append(raft.log, log)
}

func GetLog() []Log {
	return raft.log
}

func GetTerm() int {
	return raft.term
}

func SetTerm(term int) {
	raft.term = term
}

func GetRaftIsLeader() bool {
	return raft.role == Leader
}

func GetLeader() string {
	return raft.leader
}

func (r *Raft) isLeader() bool {
	return r.role == Leader
}

func StartRaft(host string, port string) {
	raft.self = host + ":" + port
	raft.resetElectionTimeout()
	go raft.Heartbeat()
}

func (r *Raft) resetElectionTimeout() {
	if r.electionTimeout != nil {
		r.electionTimeout.Stop()
	}
	r.electionTimeout = time.AfterFunc(time.Duration(5+rand.Intn(5))*time.Second, func() {
		r.mu.Lock()
		if r.role != Leader {
			r.role = Candidate
			r.term++
			r.votes = 1
			log.Printf("Node (SELF) %s became Candidate in term %d", raft.self, r.term)
			r.mu.Unlock()
			r.startElection()
		} else {
			r.mu.Unlock()
		}
	})
}

func ResetElectionTimeout() {
	raft.resetElectionTimeout()
}

func GetSelf() string {
	return raft.self
}

func GetMembers() []string {
	return raft.members
}

func SetMember(member []string) {
	raft.members = member
}

func Berlangganan(w http.ResponseWriter, r *http.Request) {
	if raft.role != Leader && r.Method == http.MethodPost {
		if raft.leader == "" {
			return
		}
		var requestURL = "http://" + raft.leader + "/subscribe"
		resp, err := http.Post(requestURL,
			"application/json",
			r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

	} else if r.Method == http.MethodPost {
		var subscribeReq SubscribeReq

		if err := json.NewDecoder(r.Body).Decode(&subscribeReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !contains(raft.members, subscribeReq.IPAddress) {
			raft.members = append(raft.members, subscribeReq.IPAddress)
		}

		var response = Response{
			Status:     "success",
			Message:    "Subscribed",
			RaftMember: raft.members,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func contains(arr []string, target string) bool {
	for _, elem := range arr {
		if elem == target {
			return true
		}
	}
	return false
}
