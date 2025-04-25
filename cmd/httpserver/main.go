package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"myhttpprotocol/internal/server"

	"github.com/pingcap/log"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
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