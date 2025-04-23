package request

import (
	"errors"
	"fmt"
	"io"
	"myhttpprotocol/internal/headers"
	"strconv"
	"strings"
)

type ParserState int

const (
	StateInitialized ParserState = iota
	StateParsingHeaders
	StateParsingBody
	StateDone
)

type Request struct {
	RequestLine RequestLine
	state       ParserState
	Headers     headers.Headers
	Body 		[]byte
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state:   StateInitialized,
		Headers: headers.NewHeaders(),
		Body:   make([]byte, 0),
	}

	buffer := make([]byte, 8)
	parsedBytes := 0

	for req.state != StateDone {
		if parsedBytes >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		n, err := reader.Read(buffer[parsedBytes:])
		if err != nil {
			if err == io.EOF {
				if req.state != StateDone {
					return nil, errors.New("incomplete request: missing end of headers")
				}
				break
			}
			return nil, err
		}

		parsedBytes += n

		consumed, err := req.parse(buffer[:parsedBytes])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[consumed:])
		parsedBytes -= consumed
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != StateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// Need more data
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case StateInitialized:
		n, requestLine, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// Need more data
			return 0, nil
		}
		r.RequestLine = requestLine
		r.state = StateParsingHeaders
		return n, nil

	case StateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = StateParsingBody
		}

		return n, nil
	
	case StateParsingBody:
		contentLength := r.Headers.Get("Content-Length")
		if contentLength == "" {
			r.state = StateDone
			return len(data), nil
		}

		num, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("invalid Content-Length value: %v", err)
		}

		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)

		if r.bodyLengthRead > num {
			return 0, fmt.Errorf("data is larger than shared Content-Length")
		}

		if r.bodyLengthRead == num {
			r.state = StateDone
		}

		return len(data), nil

	case StateDone:
		return 0, errors.New("trying to read error in completed state")

	default:
		return 0, errors.New("unknown state")
	}
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
	if !isValidMethod(method) {
		return 0, RequestLine{}, errors.New("invalid method")
	}

	target := requestParts[1]
	version := requestParts[2]

	// Validate version
	if version != "HTTP/1.1" {
		return 0, RequestLine{}, errors.New("unsupported HTTP version")
	}

	// Return bytes consumed (request line + CRLF)
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
