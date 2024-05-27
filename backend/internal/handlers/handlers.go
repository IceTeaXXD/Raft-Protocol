package handlers

import (
    "if3230-tubes2-spg/internal/store"
    rft "if3230-tubes2-spg/internal/raft"
    "fmt"
    "net/http"
)

func CheckLeader(w http.ResponseWriter) bool {
    if rft.GetRaftIsLeader() {
        return true
    } 
    var errorBody string = fmt.Sprintf("leader:%s", rft.GetLeader())
    http.Error(w, errorBody, http.StatusServiceUnavailable)
    return false
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
    if !CheckLeader(w) {
        return
    }
    if r.Method == http.MethodGet {
        fmt.Fprintf(w, "Pong")
    }
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
    if !CheckLeader(w) {
        return
    }
    key := r.URL.Query().Get("key")
    value := store.Get(key)
    fmt.Fprintf(w, "%s", value)
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
    if !CheckLeader(w) {
        return
    }
    key := r.URL.Query().Get("key")
    value := r.URL.Query().Get("value")
    store.Set(key, value)
    fmt.Fprintf(w, "OK")
}

func StrlnHandler(w http.ResponseWriter, r *http.Request) {
    if !CheckLeader(w) {
        return
    }
    key := r.URL.Query().Get("key")
    value := store.Get(key)
    fmt.Fprintf(w, "%d", len(value))
}

func DelHandler(w http.ResponseWriter, r *http.Request) {
    if !CheckLeader(w) {
        return
    }
    key := r.URL.Query().Get("key")
    value := store.Del(key)
    fmt.Fprintf(w, "%s", value)
}

func AppendHandler(w http.ResponseWriter, r *http.Request) {
    if !CheckLeader(w) {
        return
    }
    key := r.URL.Query().Get("key")
    value := r.URL.Query().Get("value")
    store.Append(key, value)
    fmt.Fprintf(w, "OK")
}
