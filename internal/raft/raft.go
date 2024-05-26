package raft

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

type Raft struct {
    members []string
    leader  string
    log     []string
    mu      sync.Mutex
}

var raft = Raft{
    members: []string{"localhost:8080"},
    leader:  "localhost:8080",
    log:     []string{},
}

func StartRaft() {
    go raft.Heartbeat()
}

func (r *Raft) AppendLog(entry string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.log = append(r.log, entry)
    fmt.Println("Log appended:", entry)
}

func (r *Raft) Heartbeat() {
    for {
        r.mu.Lock()
        for _, member := range r.members {
            if member != r.leader {
                go r.ping(member)
            }
        }
        r.mu.Unlock()
        time.Sleep(5 * time.Second)
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

// Implementasi untuk LeaderElection dan MembershipChange :v