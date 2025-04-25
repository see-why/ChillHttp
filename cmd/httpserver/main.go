package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"chillhttp/internal/request"
	"chillhttp/internal/server"

	"github.com/pingcap/log"
)

const port = 42069

func HttpHandler(w io.Writer, req *request.Request) (error *server.HandlerError) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			Code: 400,
			Err:  "Your problem is not my problem",
		}
	case "/myproblem":
		return &server.HandlerError{
			Code: 500,
			Err:  "Woopsie, my bad",
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}

func main() {
	server, err := server.Serve(port, HttpHandler)
	if err != nil {
		fmt.Println("Error starting server: %w", err)
	}
	defer server.Close()
	fmt.Println("Server listening on :", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Info("Received shutdown signal, shutting down server...")
}