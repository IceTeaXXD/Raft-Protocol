package client

import (
    "fmt"
    "io"
    "net/http"
)

func makeRequest(url string) (string, error) {
    response, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}

func Ping() (string, error) {
    return makeRequest("http://localhost:8080/ping")
}

func Get(key string) (string, error) {
    url := fmt.Sprintf("http://localhost:8080/get?key=%s", key)
    return makeRequest(url)
}

func Set(key, value string) (string, error) {
    url := fmt.Sprintf("http://localhost:8080/set?key=%s&value=%s", key, value)
    return makeRequest(url)
}

func Strln(key string) (string, error) {
    url := fmt.Sprintf("http://localhost:8080/strln?key=%s", key)
    return makeRequest(url)
}

func Del(key string) (string, error) {
    url := fmt.Sprintf("http://localhost:8080/del?key=%s", key)
    return makeRequest(url)
}

func Append(key, value string) (string, error) {
    url := fmt.Sprintf("http://localhost:8080/append?key=%s&value=%s", key, value)
    return makeRequest(url)
}
