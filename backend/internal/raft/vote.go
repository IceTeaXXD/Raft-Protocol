package raft

import (
	"bytes"
	"encoding/json"
	"if3230-tubes2-spg/internal/store"
	"log"
	"net/http"
)

type VoteRequest struct {
    Term        int
    CandidateID string
}

type VoteResponse struct {
    Term        int
    VoteGranted bool
}

func (r *Raft) requestVote(member string) bool {
    voteRequest := VoteRequest{
        Term:        r.term,
        CandidateID: r.self,
    }

    data, err := json.Marshal(voteRequest)
    if err != nil {
        return false
    }

    resp, err := http.Post("http://"+member+"/vote", "application/json", bytes.NewBuffer(data))
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    var voteResponse VoteResponse
    if err := json.NewDecoder(resp.Body).Decode(&voteResponse); err != nil {
        return false
    }

    return voteResponse.VoteGranted
}

func HandleVoteRequest(w http.ResponseWriter, req *http.Request) {
    var voteRequest VoteRequest
    if err := json.NewDecoder(req.Body).Decode(&voteRequest); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Node %s received vote request from %s", raft.self, voteRequest.CandidateID)

    raft.mu.Lock()
    defer raft.mu.Unlock()

    if raft.role == Candidate {
        log.Printf("Node %s is candidate and did not vote for %s", raft.self, voteRequest.CandidateID)
        voteResponse := VoteResponse{
            Term:        raft.term,
            VoteGranted: false,
        }
        if err := json.NewEncoder(w).Encode(voteResponse); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    raft.resetElectionTimeout()

    voteResponse := VoteResponse{
        Term: raft.term,
    }

    if voteRequest.Term > raft.term {
        raft.term = voteRequest.Term
        raft.votedFor = ""
        raft.role = Follower
    } else {
        voteResponse.VoteGranted = false
        log.Printf("Node %s did not vote for %s", raft.self, voteRequest.CandidateID)
        if err := json.NewEncoder(w).Encode(voteResponse); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        return
    }

    if raft.votedFor == "" || raft.votedFor == voteRequest.CandidateID {
        voteResponse.VoteGranted = true
        raft.votedFor = voteRequest.CandidateID
        log.Printf("Node %s voted for %s", raft.self, voteRequest.CandidateID)
    } else {
        voteResponse.VoteGranted = false
        log.Printf("Node %s did not vote for %s", raft.self, voteRequest.CandidateID)
    }

    if err := json.NewEncoder(w).Encode(voteResponse); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func HandleHeartbeat(w http.ResponseWriter, req *http.Request) {
    var heartbeat Heartbeat
    if err := json.NewDecoder(req.Body).Decode(&heartbeat); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Node %s received heartbeat from %s", raft.self, heartbeat.Sender)

    raft.leader = heartbeat.Sender
    raft.log = heartbeat.Log
    raft.members = heartbeat.Members

    store.Reset()
    
    for _, log := range raft.log {
        switch log.Command {
		case "set":
            store.Set(log.Arg1, log.Arg2)
		case "append":
            store.Append(log.Arg1, log.Arg2)
		case "get":
            store.Get(log.Arg1)
		case "strln":
            store.Strln(log.Arg1)
		case "del":
            store.Del(log.Arg1)
		}
    }

    raft.mu.Lock()
    defer raft.mu.Unlock()

    if raft.role != Leader {
        raft.resetElectionTimeout()
    }

    heartbeatResponse := HeartbeatResponse{
        Term: raft.term,
        Sender: raft.self,
    }

    if err := json.NewEncoder(w).Encode(heartbeatResponse); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
