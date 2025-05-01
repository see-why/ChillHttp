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

type Writer struct {
	Writer io.Writer
	State WriteState
}

type WriteState int

const (
	StateWriteStatusLine  WriteState = iota
	StateWriteHeaders
	StateWriteBody
	StateDone
)

// NewResponseWriter creates a new ResponseWriter instance
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer:     w,
		State:      StateWriteStatusLine,
	}
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != StateWriteBody {
		return 0, fmt.Errorf("invalid state: expected StateWriteHeaders, got %v", w.State)
	}

	length, err := w.Writer.Write([]byte(p))
	if err != nil {
		return length, err
	}

	w.State = StateDone
	return length, nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != StateWriteStatusLine {
		return fmt.Errorf("invalid state: expected StateInitialized, got %v", w.State)
	}

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

	_, err := w.Writer.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	w.State = StateWriteHeaders
	
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Type"] = "text/plain"
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != StateWriteHeaders {
		return fmt.Errorf("invalid state: expected StateWriteStatusLine, got %v", w.State)
	}
	for key, value := range headers {
		_, err := w.Writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.State = StateWriteBody
	return nil
}

// WriteChunkedBody writes a single chunk in chunked transfer encoding.
func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	// Write chunk size in hex followed by \r\n
	sizeLine := fmt.Sprintf("%x\r\n", len(p))
	if _, err := w.Writer.Write([]byte(sizeLine)); err != nil {
		return 0, err
	}
	// Write chunk data
	n, err := w.Writer.Write(p)
	if err != nil {
		return n, err
	}
	// Write trailing \r\n
	if _, err := w.Writer.Write([]byte("\r\n")); err != nil {
		return n, err
	}
	return n, nil
}

// WriteChunkedBodyDone writes the final zero-length chunk.
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Writer.Write([]byte("0\r\n\r\n"))
}
