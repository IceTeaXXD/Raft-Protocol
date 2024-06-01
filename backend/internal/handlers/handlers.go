package handlers

import (
	"fmt"
	rft "if3230-tubes2-spg/internal/raft"
	"if3230-tubes2-spg/internal/store"
	"io"
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

	if rft.GetRaftIsLeader() && r.Method == http.MethodPut {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		store.Set(key, value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response" : "success"}`))
	} else if r.Method == http.MethodPut {
		req, err := http.NewRequest(http.MethodPut, "http://"+rft.GetLeader()+"/set"+"?key="+r.URL.Query().Get("key")+"&value="+r.URL.Query().Get("value"), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		resp, err := http.DefaultClient.Do(req)
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

	if rft.GetRaftIsLeader() && r.Method == http.MethodDelete {
		key := r.URL.Query().Get("key")
		value := store.Del(key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response" : ` + value + `}`))
	} else if r.Method == http.MethodDelete {
		url := "http://" + rft.GetLeader() + "/del?key=" + r.URL.Query().Get("key")

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := client.Do(req)
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

func AppendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if rft.GetRaftIsLeader() && r.Method == http.MethodPut {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		store.Append(key, value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response" : "success"}`))
	} else if r.Method == http.MethodPut {
		url := "http://" + rft.GetLeader() + "/append?key=" + r.URL.Query().Get("key") + "&value=" + r.URL.Query().Get("value")
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		resp, err := client.Do(req)
		if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("Failed to execute request: " + err.Error()))
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
