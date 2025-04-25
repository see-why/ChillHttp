package server

import (
	"chillhttp/internal/request"
	"chillhttp/internal/response"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	Handler  Handler
	Closed  atomic.Bool
}
type HandlerError struct {
	Code int
	Err  string
}

type Handler func(w *response.Writer, req *request.Request)

func WriteError(w io.Writer, err *HandlerError) {
	if err == nil {
		return
	}

	body := err.Err
	statusCode := response.StatusCode(err.Code)
	writer := response.NewWriter(w)
	writer.WriteStatusLine(statusCode)
	writer.WriteHeaders(response.GetDefaultHeaders(len(body)))// Assuming buff.Len() is not available in this context
	w.Write([]byte(body))
}


func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error creating listener: %w", err)
	}

	s := &Server{
		Listener: l,
		Handler:  handler,
		Closed:  atomic.Bool{},
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	if s.Listener != nil {
		s.Closed.Store(true)
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	if s.Closed.Load() {
		return
	}
	defer conn.Close()

	writer := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.BadRequest)
		writer.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}

	writer = response.NewWriter(conn)
	s.Handler(writer, req)
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