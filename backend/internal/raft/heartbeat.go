package raft

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Heartbeat struct {
    Term int
    Sender string
    Log []Log
}

type HeartbeatResponse struct {
    Term int
    Sender string
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
        Sender: r.self,
        Log: r.log,
    }
    // Marshalling failss here with r.log
    data, err := json.Marshal(heartbeat)
    if err != nil {
        log.Printf("Failed to marshal heartbeat: %v", err)
        return
    }

    resp, err := http.Post("http://"+member+"/heartbeat", "application/json", bytes.NewBuffer(data))
    if err != nil || resp.StatusCode != http.StatusOK {
        log.Printf("Failed to send heartbeat to %s: %v", member, err)

        r.mu.Lock()
        for i, m := range r.members {
            if m == member {
                r.members = append(r.members[:i], r.members[i+1:]...)
                break
            }
        }
        r.mu.Unlock()
    }
}
