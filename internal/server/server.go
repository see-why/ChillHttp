package server

import (
	"bytes"
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
		Closed:  false,
	}
	go s.listen(handler)
	return s, nil
}

func (s *Server) Close() error {
	if s.Listener != nil {
		s.Closed = true
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	if s.Closed {
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
	herr := handler(&buff, req)
	if herr != nil {
		WriteError(conn, herr)
		return
	}

	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(buff.Len()))
	conn.Write(buff.Bytes())
}

func (s *Server) listen(handler Handler) {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: %w", err)
		}

		go s.handle(conn, handler)
	}
}