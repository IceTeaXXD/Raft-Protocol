package raft

import (
    "log"
    "sync"
)

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
                    if(!contains(r.follower, member)){
                        r.follower = append(r.follower, member)
                    }
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
        log.Printf("Node (SELF) %s became Leader in term %d\n", r.self, r.term)
        log.Printf("Members: %v\n", r.members)
    } else {
        r.role = Follower
        log.Println("Election failed, staying as candidate")
        r.resetElectionTimeout()
    }
    r.mu.Unlock()
}
