package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

const port = 42069

type Server struct {
	Listener net.Listener
	Closed  bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %w", err)
	}

	s := &Server{
		Listener: l,
		Closed:  false,
	}

	return s, nil
}

func (s *Server) Close() error {
	if s.Listener != nil {
		s.Closed = true
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	if s.Closed {
		return
	}
	defer conn.Close()

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/plain \r\n"))
	conn.Write([]byte("\nHello World!"))
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: %w", err)
		}

		go s.handle(conn)
	}
}

func main() {
	server, err := Serve(port)
	if err != nil {
		fmt.Println("Error starting server: %w", err)
	}
	defer server.Close()
	log.Info("Server listening on :%d", zap.Int("port", port))

	server.listen()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Info("Received shutdown signal, shutting down server...")
}