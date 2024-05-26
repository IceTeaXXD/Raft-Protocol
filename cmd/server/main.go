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
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file")
    }

    // Set up HTTP server
    http.HandleFunc("/ping", handlers.PingHandler)
    http.HandleFunc("/get", handlers.GetHandler)
    http.HandleFunc("/set", handlers.SetHandler)
    http.HandleFunc("/strln", handlers.StrlnHandler)
    http.HandleFunc("/del", handlers.DelHandler)
    http.HandleFunc("/append", handlers.AppendHandler)

    // Start Raft consensus
    go raft.StartRaft()

    // Get port from environment variable
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not set
    }
    log.Println("Starting server on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}