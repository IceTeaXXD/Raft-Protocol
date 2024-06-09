package handlers

import (
	"encoding/json"
	"fmt"
	rft "if3230-tubes2-spg/internal/raft"
	"if3230-tubes2-spg/internal/store"
	"io"
	"log"
	"net/http"
	"strconv"
	// "strings"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if rft.GetRaftIsLeader() && r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "Pong"}`))
		rft.AddLog("ping", "", "")
	} else if r.Method == http.MethodGet {
		resp, err := http.Get("http://" + rft.GetLeader() + "/ping")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method"))
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if rft.GetRaftIsLeader() && r.Method == http.MethodGet {
		key := r.URL.Query().Get("key")
		value := store.Get(key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte((`{"response": "` + value + `"}`)))
		rft.AddLog("get", key, "")
	} else if r.Method == http.MethodGet {
		resp, err := http.Get("http://" + rft.GetLeader() + "/get" + "?key=" + r.URL.Query().Get("key"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method"))
	}
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key and value are required"}`))
		return
	}

	if rft.GetRaftIsLeader() {
		term := rft.GetTerm()
		// Send request to all members
		members := rft.GetMembers()
		for _, m := range members {
			if m == rft.GetLeader() {
				continue
			}
			req, err := http.NewRequest(http.MethodPut, "http://"+m+"/setReplicate"+"?key="+key+"&value="+value+"&term="+strconv.Itoa(term), nil)
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Failed to replicate to member %s: %v", m, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
				return
			}
			resp.Body.Close()
		}

		store.Set(key, value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "success"}`))
		rft.AddLog("set", key, value)
	} else {
		leader := rft.GetLeader()
		req, err := http.NewRequest(http.MethodPut, "http://"+leader+"/set"+"?key="+key+"&value="+value, nil)
		if err != nil {
			log.Printf("Failed to create request to leader: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to forward request to leader: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}

func replicateSetToFollowers(key, value string, term int) {
	members := rft.GetMembers()
	for _, m := range members {
		if m == rft.GetLeader() {
			continue
		}
		req, err := http.NewRequest(http.MethodPut, "http://"+m+"/setReplicate"+"?key="+key+"&value="+value+"&term="+strconv.Itoa(term), nil)
		if err != nil {
			log.Printf("Failed to create replication request to member %s: %v", m, err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to replicate set to member %s: %v", m, err)
			continue
		}
		resp.Body.Close()
	}
}

func SetReplicateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	term := r.URL.Query().Get("term")

	if key == "" || value == "" || term == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key, value, and term are required"}`))
		return
	}

	termInt, err := strconv.Atoi(term)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid term"}`))
		return
	}

	// Ngecek term dulu, kalau termnya lebih kecil dari term sekarang ya berarti ketinggalan jaman :V
	currentTerm := rft.GetTerm()
	if termInt < currentTerm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Outdated term"}`))
		return
	}

	if termInt > currentTerm {
		rft.SetTerm(termInt)
	}

	store.Set(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "success"}`))
}

func StrlnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if rft.GetRaftIsLeader() && r.Method == http.MethodGet {
		key := r.URL.Query().Get("key")
		value := store.Get(key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": ` + fmt.Sprint(len(value)) + `}`))
		rft.AddLog("strln", key, "")
	} else if r.Method == http.MethodGet {
		resp, err := http.Get("http://" + rft.GetLeader() + "/strln" + "?key=" + r.URL.Query().Get("key"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method"))
	}
}

func DelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key is required"}`))
		return
	}

	if rft.GetRaftIsLeader() {
		term := rft.GetTerm()
		value := store.Del(key)
		replicateDelToFollowers(key, term)
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"response": value}
		jsonResp, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		w.Write(jsonResp)
		rft.AddLog("del", key, "")
	} else {
		url := "http://" + rft.GetLeader() + "/del?key=" + key

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to execute request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}

func replicateDelToFollowers(key string, term int) {
	members := rft.GetMembers()
	for _, m := range members {
		if m == rft.GetLeader() {
			continue
		}
		req, err := http.NewRequest(http.MethodDelete, "http://"+m+"/delReplicate"+"?key="+key+"&term="+strconv.Itoa(term), nil)
		if err != nil {
			log.Printf("Failed to create replication request to member %s: %v", m, err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to replicate delete to member %s: %v", m, err)
			continue
		}
		resp.Body.Close()
	}
}

func DelReplicateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	key := r.URL.Query().Get("key")
	term := r.URL.Query().Get("term")

	if key == "" || term == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key and term are required"}`))
		return
	}

	termInt, err := strconv.Atoi(term)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid term"}`))
		return
	}

	// Ngecek term dulu, kalau termnya lebih kecil dari term sekarang ya berarti ketinggalan jaman :V
	currentTerm := rft.GetTerm()
	if termInt < currentTerm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Outdated term"}`))
		return
	}

	if termInt > currentTerm {
		rft.SetTerm(termInt)
	}

	value := store.Del(key)
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"response": value}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}
	w.Write(jsonResp)
}
func AppendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key and value are required"}`))
		return
	}

	if rft.GetRaftIsLeader() {
		term := rft.GetTerm()
		store.Append(key, value)
		replicateAppendToFollowers(key, value, term)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "success"}`))
		rft.AddLog("append", key, value)
	} else {
		url := "http://" + rft.GetLeader() + "/append?key=" + key + "&value=" + value
		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to execute request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}

func replicateAppendToFollowers(key, value string, term int) {
	members := rft.GetMembers()
	for _, m := range members {
		if m == rft.GetLeader() {
			continue
		}
		req, err := http.NewRequest(http.MethodPut, "http://"+m+"/appendReplicate"+"?key="+key+"&value="+value+"&term="+string(term), nil)
		if err != nil {
			log.Printf("Failed to create replication request to member %s: %v", m, err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to replicate log to member %s: %v", m, err)
			continue
		}
		resp.Body.Close()
	}
}

func AppendReplicateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	term := r.URL.Query().Get("term")

	if key == "" || value == "" || term == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key, value, and term are required"}`))
		return
	}

	termInt, err := strconv.Atoi(term)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid term"}`))
		return
	}

	// Ngecek term dulu, kalau termnya lebih kecil dari term sekarang ya berarti ketinggalan jaman :V
	currentTerm := rft.GetTerm()
	if termInt < currentTerm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Outdated term"}`))
		return
	}

	if termInt > currentTerm {
		rft.SetTerm(termInt)
	}

	store.Append(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "success"}`))
}

func RequestLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid method"}`))
		return
	}

	if rft.GetRaftIsLeader() {
		w.WriteHeader(http.StatusOK)
		// w.Write([]byte(`{"response": "` + strings.Join(rft.GetLog(), " | ") + `"}`))
	} else {
		// Forward request to leader
		url := "http://" + rft.GetLeader() + "/requestLog"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to execute request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}
