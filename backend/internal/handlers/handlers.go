package handlers

import (
	"encoding/json"
	"fmt"
	rft "if3230-tubes2-spg/internal/raft"
	"if3230-tubes2-spg/internal/store"
	"io"
	"log"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if rft.GetRaftIsLeader() && r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "Pong"}`))
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
		// Send request to all members
		members := rft.GetMembers()
		for _, m := range members {
			if m == rft.GetLeader() {
				continue
			}
			req, err := http.NewRequest(http.MethodPut, "http://"+m+"/setReplicate"+"?key="+key+"&value="+value, nil)
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

func SetReplicateHandler(w http.ResponseWriter, r *http.Request) {
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
		value := store.Del(key)
		replicateDelToFollowers(key)
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

func replicateDelToFollowers(key string) {
	members := rft.GetMembers()
	for _, m := range members {
		if m == rft.GetLeader() {
			continue
		}
		req, err := http.NewRequest(http.MethodDelete, "http://"+m+"/delReplicate"+"?key="+key, nil)
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

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key is required"}`))
		return
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
		store.Append(key, value)
		replicateAppendToFollowers(key, value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "success"}`))
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

func replicateAppendToFollowers(key, value string) {
	members := rft.GetMembers()
	for _, m := range members {
		if m == rft.GetLeader() {
			continue
		}
		req, err := http.NewRequest(http.MethodPut, "http://"+m+"/appendReplicate"+"?key="+key+"&value="+value, nil)
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

	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Key and value are required"}`))
		return
	}

	store.Append(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "success"}`))
}