package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"myhttpprotocol/internal/server"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
	if err != nil {
		fmt.Println("Error starting server: %w", err)
	}
	defer server.Close()
	log.Info("Server listening on :%d", zap.Int("port", port))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Info("Received shutdown signal, shutting down server...")
}