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
                if member != r.leader {
                    go r.sendHeartbeat(member)
                }
            }
        }
        r.mu.Unlock()
    }
}

func (r *Raft) sendHeartbeat(member string) {
    resp, err := http.Get("http://" + member + "/ping")
    if err != nil || resp.StatusCode != http.StatusOK {
        log.Printf("Failed to ping %s: %v", member, err)
        return
    }
    r.resetElectionTimeout()
}

func (r *Raft) startElection() {
    var wg sync.WaitGroup

    for _, member := range r.members {
        if member != r.self {
            wg.Add(1)
            go func(member string) {
                defer wg.Done()
                log.Printf("Node %s requesting vote from %s", "self", member)
                if r.requestVote(member) {
                    r.mu.Lock()
                    r.votes++
                    log.Printf("Node %s received vote from %s", "self", member)
                    r.mu.Unlock()
                }
            }(member)
        }
    }

    wg.Wait()

    r.mu.Lock()
    if r.votes > len(r.members)/2 {
        r.role = Leader
        r.leader = "self"
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
        CandidateID: "self",
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

    raft.mu.Lock()
    defer raft.mu.Unlock()

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
    } else {
        voteResponse.VoteGranted = false
    }

    if err := json.NewEncoder(w).Encode(voteResponse); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
