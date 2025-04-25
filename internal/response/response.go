package response

import (
	"chillhttp/internal/headers"
	"fmt"
	"io"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
		case OK:
			statusLine = "HTTP/1.1 200 OK\r\n"
		case BadRequest:
			statusLine = "HTTP/1.1 400 Bad Request\r\n"
		case InternalServerError:
			statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
		default:
			statusLine = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Type"] = "text/plain"
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}