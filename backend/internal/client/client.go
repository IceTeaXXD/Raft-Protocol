package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	Host          string
	Port          string
	InTransaction bool
	Transaction   []Command
}

type Command struct {
	Name string
	Args []string
}

type JsonResponse struct {
	Response json.RawMessage `json:"response"`
}

func (c *Client) makeRequest(method, endpoint string, body io.Reader) (string, error) {
	url := "http://" + c.Host + ":" + c.Port + endpoint
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

	fmt.Println(responseStr)

	var responseNumber int64
	if err := json.Unmarshal(jsonResponse.Response, &responseNumber); err == nil {
		return strconv.FormatInt(responseNumber, 10), nil
	}

	return "", err
}

func (c *Client) Ping() (string, error) {
	if c.InTransaction {
		c.Transaction = append(c.Transaction, Command{Name: "Ping", Args: []string{}})
		return "", nil
	}
	return c.makeRequest(http.MethodGet, "/ping", nil)
}

func (c *Client) Get(key string) (string, error) {
	if c.InTransaction {
		c.Transaction = append(c.Transaction, Command{Name: "Get", Args: []string{key}})
		return "", nil
	}
	url := fmt.Sprintf("/get?key=%s", key)
	return c.makeRequest(http.MethodGet, url, nil)
}

func (c *Client) Set(key, value string) (string, error) {
	if c.InTransaction {
		c.Transaction = append(c.Transaction, Command{Name: "Set", Args: []string{key, value}})
		return "", nil
	}
	url := fmt.Sprintf("/set?key=%s&value=%s", key, value)
	return c.makeRequest(http.MethodPut, url, nil)
}

func (c *Client) Strln(key string) (string, error) {
	if c.InTransaction {
		c.Transaction = append(c.Transaction, Command{Name: "Strln", Args: []string{key}})
		return "", nil
	}
	url := fmt.Sprintf("/strln?key=%s", key)
	return c.makeRequest(http.MethodGet, url, nil)
}

func (c *Client) Del(key string) (string, error) {
	if c.InTransaction {
		c.Transaction = append(c.Transaction, Command{Name: "Del", Args: []string{key}})
		return "", nil
	}
	url := fmt.Sprintf("/del?key=%s", key)
	return c.makeRequest(http.MethodDelete, url, nil)
}

func (c *Client) Append(key, value string) (string, error) {
	if c.InTransaction {
		c.Transaction = append(c.Transaction, Command{Name: "Append", Args: []string{key, value}})
		return "", nil
	}
	url := fmt.Sprintf("/append?key=%s&value=%s", key, value)
	return c.makeRequest(http.MethodPut, url, nil)
}

func (c* Client) RequestLog() (string, error) {
	return c.makeRequest(http.MethodGet, "/requestLog", nil)
}
func (c *Client) Begin() {
	c.InTransaction = true
}

func (c *Client) Commit() ([]string, error) {
	var results []string
	c.InTransaction = false
	for _, command := range c.Transaction {
		var res string
		var err error
		switch command.Name {
		case "Set":
			res, err = c.Set(command.Args[0], command.Args[1])
		case "Append":
			res, err = c.Append(command.Args[0], command.Args[1])
		case "Get":
			res, err = c.Get(command.Args[0])
		case "Strln":
			res, err = c.Strln(command.Args[0])
		case "Del":
			res, err = c.Del(command.Args[0])
		case "Ping":
			res, err = c.Ping()
		}
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	c.Transaction = []Command{}
	return results, nil
}
