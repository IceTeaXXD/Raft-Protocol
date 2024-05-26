package main

import (
    "if3230-tubes2-spg/internal/handlers"
    "if3230-tubes2-spg/internal/raft"
    "log"
    "net/http"
    "os"
    
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    _ = godotenv.Load()

    // Set up HTTP server
    http.HandleFunc("/ping", handlers.PingHandler)
    http.HandleFunc("/get", handlers.GetHandler)
    http.HandleFunc("/set", handlers.SetHandler)
    http.HandleFunc("/strln", handlers.StrlnHandler)
    http.HandleFunc("/del", handlers.DelHandler)
    http.HandleFunc("/append", handlers.AppendHandler)

    // Ini buat Raft
    http.HandleFunc("/vote", raft.HandleVoteRequest)

    
    // Get port from environment variable
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not set
    }

    // Start Raft consensus
    go raft.StartRaft(port)
    
    log.Println("Starting server on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}