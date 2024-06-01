package client

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strings"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func getServerURL() string {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}

func makeRequest(endpoint string) (string, error) {
	url := getServerURL() + endpoint
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		leaderPort := strings.TrimSpace(string(body)[len("leader:localhost:"):])
		url = fmt.Sprintf("http://localhost:%s%s", leaderPort, endpoint)
		response, err := http.Get(url)
		if err != nil {
			return "", err
		}
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		defer response.Body.Close()
	}
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Ping() (string, error) {
	return makeRequest("/ping")
}

func Get(key string) (string, error) {
	url := fmt.Sprintf("/get?key=%s", key)
	return makeRequest(url)
}

func Set(key, value string) (string, error) {
	url := fmt.Sprintf("/set?key=%s&value=%s", key, value)
	return makeRequest(url)
}

func Strln(key string) (string, error) {
	url := fmt.Sprintf("/strln?key=%s", key)
	return makeRequest(url)
}

func Del(key string) (string, error) {
	url := fmt.Sprintf("/del?key=%s", key)
	return makeRequest(url)
}

func Append(key, value string) (string, error) {
	url := fmt.Sprintf("/append?key=%s&value=%s", key, value)
	return makeRequest(url)
}
