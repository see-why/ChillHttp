package main

import (
	"chillhttp/internal/request"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error creating listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Println("Connection accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error reading request: %v\n", err)
			conn.Close()
			continue
		}

		fmt.Printf("Request line:\n"+
			"- Method: %s\n"+
			"- Target: %s\n"+
			"- Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		if len(req.Body) > 0 {
			fmt.Println(string(req.Body))
		} else {
			fmt.Println("No body")
		}
		fmt.Println("Request processing complete")

		conn.Close()
		fmt.Println("Connection closed")
	}
}
