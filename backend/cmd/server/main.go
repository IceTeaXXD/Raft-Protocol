package main

import (
	"flag"
	"if3230-tubes2-spg/internal/handlers"
	"if3230-tubes2-spg/internal/raft"
	"log"
	"net/http"
)

func main() {
	// Parse flags
	var port string
	flag.StringVar(&port, "port", "8080", "Port to run the server on")
	flag.Parse()

	// Set up HTTP server
	http.HandleFunc("/ping", handlers.PingHandler)
	http.HandleFunc("/get", handlers.GetHandler)
	http.HandleFunc("/set", handlers.SetHandler)
	http.HandleFunc("/strln", handlers.StrlnHandler)
	http.HandleFunc("/del", handlers.DelHandler)
	http.HandleFunc("/append", handlers.AppendHandler)

	// Ini buat Raft
	http.HandleFunc("/vote", raft.HandleVoteRequest)
	http.HandleFunc("/heartbeat", raft.HandleHeartbeat)

	// Start Raft consensus
	go raft.StartRaft(port)

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
