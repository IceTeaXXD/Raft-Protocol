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

func main() {
	// Parse flags
	var port string
	var lead string
	var portLead string

	flag.StringVar(&port, "port", "8080", "Port to run the server on")
	flag.StringVar(&lead, "lead", "0.0.0.0", "Port to be the leader of this server")
	flag.StringVar(&portLead, "portLead", "8080", "Lead port")
	flag.Parse()

	if (port != "8080"){
		var subscribeReq = raft.SubscribeReq{
			IPAddress: lead + ":" + port,
		}
		
		payload, _ := json.Marshal(subscribeReq)
		
		var url = "http://" + lead  + ":" + portLead + "/subscribe"
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

	// Start Raft consensus
	go raft.StartRaft(port)

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
