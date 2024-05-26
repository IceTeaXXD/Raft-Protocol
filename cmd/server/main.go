package main

import (
    "if3230-tubes2-spg/internal/handlers"
    "if3230-tubes2-spg/internal/raft"
    "log"
    "net/http"
)

func main() {
    // Set up HTTP server
    http.HandleFunc("/ping", handlers.PingHandler)
    http.HandleFunc("/get", handlers.GetHandler)
    http.HandleFunc("/set", handlers.SetHandler)
    http.HandleFunc("/strln", handlers.StrlnHandler)
    http.HandleFunc("/del", handlers.DelHandler)
    http.HandleFunc("/append", handlers.AppendHandler)

    // Start Raft consenSUS
    go raft.StartRaft()

    // Start HTTP server
    port := ":8080"
    log.Println("Starting server on port", port)
    log.Fatal(http.ListenAndServe(port, nil))
}