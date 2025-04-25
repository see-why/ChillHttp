package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteError(w io.Writer, err *HandlerError) {
	if err == nil {
		return
	}

	body := err.Err
	statusCode := response.StatusCode(err.Code)
	response.WriteStatusLine(w, statusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(body)))// Assuming buff.Len() is not available in this context
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.BadRequest)
		response.WriteHeaders(conn, response.GetDefaultHeaders(0))
		return
	}

	var buff bytes.Buffer
	herr := s.Handler(&buff, req)
	if herr != nil {
		WriteError(conn, herr)
		return
	}

	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(buff.Len()))
	conn.Write(buff.Bytes())
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