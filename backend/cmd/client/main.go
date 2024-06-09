package main

import (
	"bufio"
	"flag"
	"fmt"
	"if3230-tubes2-spg/internal/client"
	"os"
	"strings"
)

func main() {
	// Parse the port from the command line arguments
	var port string
	var host string
	flag.StringVar(&port, "port", "8080", "Port to run the client on")
	flag.StringVar(&host, "host", "localhost", "Host Server")
	flag.Parse()

	// Create a new client instance with the specified port
	c := &client.Client{Host: host, Port: port}

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Client started. Enter commands:")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		var res string
		var err error

		switch cmd {
		case "set":
			if len(args) != 2 {
				fmt.Println("Usage: set <key> <value>")
				continue
			}
			res, err = c.Set(args[0], args[1])
		case "append":
			if len(args) != 2 {
				fmt.Println("Usage: append <key> <value>")
				continue
			}
			res, err = c.Append(args[0], args[1])
		case "get":
			if len(args) != 1 {
				fmt.Println("Usage: get <key>")
				continue
			}
			res, err = c.Get(args[0])
		case "strln":
			if len(args) != 1 {
				fmt.Println("Usage: strln <key>")
				continue
			}
			res, err = c.Strln(args[0])
		case "del":
			if len(args) != 1 {
				fmt.Println("Usage: del <key>")
				continue
			}
			res, err = c.Del(args[0])
		case "ping":
			res, err = c.Ping()
		case "log":
			res, err = c.RequestLog()
		case "exit":
			fmt.Println("Exiting client.")
			return
		default:
			fmt.Println("Unknown command:", cmd)
			continue
		}

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(res)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
