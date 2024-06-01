package raft

import (
    "log"
    "math/rand"
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
    members         []string
    follower        []string
    leader          string
    self            string
    log             []string
    mu              sync.Mutex
    role            NodeRole
    term            int
    votedFor        string
    votes           int
    electionTimeout *time.Timer
}

var raft = Raft{
    members:  []string{"localhost:8081", "localhost:8082", "localhost:8083"},
    leader:   "",
    log:      []string{},
    role:     Follower,
    term:     0,
    votedFor: "",
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

func ResetElectionTimeout() {
    raft.resetElectionTimeout()
}

func GetSelf() string {
    return raft.self
}
