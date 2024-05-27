package raft

import (
    "bytes"
    "encoding/json"
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

type VoteRequest struct {
    Term        int
    CandidateID string
}

type VoteResponse struct {
    Term        int
    VoteGranted bool
}

type Heartbeat struct {
    Term int
}

type HeartbeatResponse struct {
    Term int
}

type Raft struct {
    members        []string
    leader         string
    self           string
    log            []string
    mu             sync.Mutex
    role           NodeRole
    term           int
    votedFor       string
    votes          int
    electionTimeout *time.Timer
}

var raft = Raft{
    members:  []string{"localhost:8080", "localhost:8081", "localhost:8082"},
    leader:   "",
    log:      []string{},
    role:     Follower,
    term:     0,
    votedFor: "",
}

func StartRaft(port string) {
    raft.self = "localhost:" + port
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

func (r *Raft) Heartbeat() {
    for {
        time.Sleep(2 * time.Second)
        r.mu.Lock()
        if r.role == Leader {
            for _, member := range r.members {
                if member != r.self {
                    go r.sendHeartbeat(member)
                }
            }
        }
        r.mu.Unlock()
    }
}

func (r *Raft) sendHeartbeat(member string) {
    log.Printf("Node %s sending heartbeat to %s", r.self, member)
    heartbeat := Heartbeat{
        Term: r.term,
    }
    data, err := json.Marshal(heartbeat)
    if err != nil {
        log.Printf("Failed to marshal heartbeat: %v", err)
        return
    }

    resp, err := http.Post("http://"+member+"/heartbeat", "application/json", bytes.NewBuffer(data))
    if err != nil || resp.StatusCode != http.StatusOK {
        log.Printf("Failed to send heartbeat to %s: %v", member, err)
        return
    }
}

func (r *Raft) startElection() {
    var wg sync.WaitGroup

    for _, member := range r.members {
        if member != r.self {
            wg.Add(1)
            go func(member string) {
                defer wg.Done()
                log.Printf("Node %s requesting vote from %s", r.self, member)
                if r.requestVote(member) {
                    r.mu.Lock()
                    r.votes++
                    log.Printf("Node %s received vote from %s", r.self, member)
                    r.mu.Unlock()
                }
            }(member)
        }
    }

    wg.Wait()

    r.mu.Lock()
    if r.votes > len(r.members)/2 {
        r.role = Leader
        r.leader = r.self
        log.Println("Became leader")
    } else {
        r.role = Follower
        log.Println("Election failed, staying as candidate")
        r.resetElectionTimeout()
    }
    r.mu.Unlock()
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

    // Kala leader, langsung deny
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

    // Reset timeout karena ada request vote masuk
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

    raft.mu.Lock()
    defer raft.mu.Unlock()

    // Kalau dapet heartbeat dari leader, reset election timeout
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

func HandlePing(w http.ResponseWriter, req *http.Request) {
    log.Printf("Node %s received ping from %s", raft.self, req.RemoteAddr)
    raft.resetElectionTimeout()
    w.WriteHeader(http.StatusOK)
}
