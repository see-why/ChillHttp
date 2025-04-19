package request

import (
	"errors"
	"io"
	"strings"
)

type ParserState int

const (
	StateInitialized ParserState = iota
	StateDone
)

type Request struct {
	RequestLine RequestLine
	state       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state: StateInitialized,
	}

	buffer := make([]byte, 8)
	parsedBytes := 0

	for req.state != StateDone {
		n, err := reader.Read(buffer[parsedBytes:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		consumed, err := req.parse(buffer[:parsedBytes+n])
		if err != nil {
			return nil, err
		}

		// Shift unparsed data to the beginning of the buffer
		if consumed < len(buffer) {
			copy(buffer, buffer[consumed:])
			parsedBytes = len(buffer) - consumed
		} else {
			parsedBytes = 0
		}

		// Grow buffer if needed
		if parsedBytes == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == StateDone {
		return 0, nil
	}

	consumed, requestLine, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if consumed == 0 {
		return 0, nil
	}

	r.RequestLine = requestLine
	r.state = StateDone
	return consumed, nil
}

func parseRequestLine(data []byte) (int, RequestLine, error) {
	// Find the end of the request line
	end := strings.Index(string(data), "\r\n")
	if end == -1 {
		return 0, RequestLine{}, nil
	}

	// Split request line into components
	requestParts := strings.Split(string(data[:end]), " ")
	if len(requestParts) != 3 {
		return 0, RequestLine{}, errors.New("invalid request line format")
	}

	method := requestParts[0]
	target := requestParts[1]
	version := requestParts[2]

	// Validate method
	if !isValidMethod(method) {
		return 0, RequestLine{}, errors.New("invalid method")
	}

	// Validate version
	if version != "HTTP/1.1" {
		return 0, RequestLine{}, errors.New("unsupported HTTP version")
	}

	// + 2 to account for \r\n characters
	return end + 2, RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}

func isValidMethod(method string) bool {
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}
