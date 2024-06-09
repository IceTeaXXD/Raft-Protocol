package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"if3230-tubes2-spg/internal/handlers"
	"if3230-tubes2-spg/internal/raft"
	"io"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	var port string
	var host string
	var leaderHost string
	var leaderPort string

	flag.StringVar(&port, "port", "8080", "Port to run the server on")
	flag.StringVar(&host, "host", "localhost", "Server Host")
	flag.StringVar(&leaderHost, "leaderHost", "", "Host of the server leader")
	flag.StringVar(&leaderPort, "leaderPort", "8080", "Port of the server leader")
	flag.Parse()

	if (leaderHost != ""){
		var subscribeReq = raft.SubscribeReq{
			IPAddress: host + ":" + port,
		}

		payload, _ := json.Marshal(subscribeReq)

		var url = "http://" + leaderHost  + ":" + leaderPort + "/subscribe"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error creating new request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if(err != nil){
			fmt.Println("Error making the request:", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var responseJSON raft.Response
		err = json.Unmarshal(body, &responseJSON)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		raftMember := responseJSON.RaftMember

		raft.SetMember(raftMember)
	}

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

	// Handle function
	http.HandleFunc("/subscribe", raft.Berlangganan)

	// Request Log Endpoint
	http.HandleFunc("/requestLog", handlers.RequestLog)

	// Log Replication
	http.HandleFunc("/setReplicate", handlers.SetReplicateHandler)
	http.HandleFunc("/delReplicate", handlers.DelReplicateHandler)
	http.HandleFunc("/appendReplicate", handlers.AppendReplicateHandler)
	// Start Raft consensus
	go raft.StartRaft(host, port)

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, enableCORS(http.DefaultServeMux)))
}
