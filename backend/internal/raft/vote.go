package raft

import (
    "bytes"
    "encoding/json"
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

    // Bruhhh, kalo leader ngapain vote
    if raft.role == Leader {
        log.Printf("Node %s is leader and did not vote for %s", raft.self, voteRequest.CandidateID)
        voteResponse := VoteResponse{
            Term:        raft.term,
            VoteGranted: false,
        }
        if err := json.NewEncoder(w).Encode(voteResponse); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Dapet vote request dari candidate, jadi reset election timeout
    raft.resetElectionTimeout()

    voteResponse := VoteResponse{
        Term: raft.term,
    }

    if voteRequest.Term > raft.term {
        raft.term = voteRequest.Term
        raft.votedFor = ""
        raft.role = Follower
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

    log.Printf("Node %s received heartbeat from %s", raft.self, req.RemoteAddr)
    // TODO: leader tuh harusnya di set sesuai heartbeat dapet dari mana 
    // raft.leader = req.TODO

    raft.mu.Lock()
    defer raft.mu.Unlock()

    // Dapet heartbeat dari leader, jadi reset election timeout
    if raft.role != Leader {
        raft.resetElectionTimeout()
    }

    heartbeatResponse := HeartbeatResponse{
        Term: raft.term,
    }

    if err := json.NewEncoder(w).Encode(heartbeatResponse); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
