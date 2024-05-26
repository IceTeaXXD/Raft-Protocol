package raft

import (
	"fmt"
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

type Raft struct {
    members []string
    leader  string
    log     []string
    mu      sync.Mutex
    role    NodeRole
    term    int
    votedFor string
}

var raft = Raft{
    members: []string{"localhost:8080", "localhost:8081", "localhost:8082"},
    leader:  "",
    log:     []string{},
    role:    Follower,
    term:    0,
    votedFor: "",
}

func StartRaft() {
    go raft.Heartbeat()
    go raft.StartElectionTimeout()
}

func (r *Raft) AppendLog(entry string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.log = append(r.log, entry)
    fmt.Println("Log appended:", entry)
}

func (r *Raft) Heartbeat() {
    for {
        if r.role == Leader {
            r.mu.Lock()
            for _, member := range r.members {
                if member != r.leader {
                    go r.ping(member)
                }
            }
            r.mu.Unlock()
        }
        time.Sleep(2 * time.Second)
    }
}

func (r *Raft) ping(member string) {
    resp, err := http.Get("http://" + member + "/ping")
    if err != nil {
        log.Printf("Failed to ping %s: %v", member, err)
        return
    }
    if resp.StatusCode == http.StatusOK {
        fmt.Println(member, "is alive")
    }
}

func (r *Raft) StartElectionTimeout() {
    for {
        time.Sleep(time.Duration(5+rand.Intn(5)) * time.Second)
        if r.role != Leader {
            r.StartElection()
        }
    }
}

func (r *Raft) StartElection() {
    r.mu.Lock()
    r.role = Candidate
    r.term++
    r.votedFor = "self"
    r.mu.Unlock()

    votes := 1
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, member := range r.members {
        if member != r.leader {
            wg.Add(1)
            go func(member string) {
                defer wg.Done()
                if r.requestVote(member) {
                    mu.Lock()
                    votes++
                    mu.Unlock()
                }
            }(member)
        }
    }

    wg.Wait()

    if votes > len(r.members)/2 {
        r.mu.Lock()
        r.role = Leader
        r.leader = "self"
        r.mu.Unlock()
        fmt.Println("Became leader")
    } else {
        r.mu.Lock()
        r.role = Follower
        r.mu.Unlock()
    }
}

func (r *Raft) requestVote(member string) bool {
    // Implement request vote RPC call
    // Vote auto terima dulu ya ges yak
    return true
}
