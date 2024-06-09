package client

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"strconv"
)

type Client struct {
	Host string
	Port string
}

type JsonResponse struct {
	Response json.RawMessage `json:"response"`
}

func (c *Client) makeRequest(method, endpoint string, body io.Reader) (string, error) {
	url := "http://"+ c.Host + ":" + c.Port + endpoint
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var jsonResponse JsonResponse
	err = json.Unmarshal(bodyBytes, &jsonResponse)
	if err != nil {
		return "", err
	}

	var responseStr string
	if err := json.Unmarshal(jsonResponse.Response, &responseStr); err == nil {
		return responseStr, nil
	}

	var responseNumber int64
	if err := json.Unmarshal(jsonResponse.Response, &responseNumber); err == nil {
		return strconv.FormatInt(responseNumber, 10), nil
	}

	return "", err
}

func (c *Client) Ping() (string, error) {
	return c.makeRequest(http.MethodGet, "/ping", nil)
}

func (c *Client) Get(key string) (string, error) {
	url := fmt.Sprintf("/get?key=%s", key)
	return c.makeRequest(http.MethodGet, url, nil)
}

func (c *Client) Set(key, value string) (string, error) {
	url := fmt.Sprintf("/set?key=%s&value=%s", key, value)
	return c.makeRequest(http.MethodPut, url, nil)
}

func (c *Client) Strln(key string) (string, error) {
	url := fmt.Sprintf("/strln?key=%s", key)
	return c.makeRequest(http.MethodGet, url, nil)
}

func (c *Client) Del(key string) (string, error) {
	url := fmt.Sprintf("/del?key=%s", key)
	return c.makeRequest(http.MethodDelete, url, nil)
}

func (c *Client) Append(key, value string) (string, error) {
	url := fmt.Sprintf("/append?key=%s&value=%s", key, value)
	return c.makeRequest(http.MethodPut, url, nil)
}
