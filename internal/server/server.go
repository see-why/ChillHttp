package server

import (
	"chillhttp/internal/request"
	"chillhttp/internal/response"
	"fmt"
	"io"
	"net"
)

type Server struct {
	Listener net.Listener
	Closed  bool
}
type HandlerError struct {
	Code int
	Err  error
}

type Handler func(w io.Writer, req *request.Request) *HandlerError


func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %w", err)
	}

	s := &Server{
		Listener: l,
		Closed:  false,
	}
	go s.listen()
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

	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
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