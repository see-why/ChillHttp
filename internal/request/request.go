package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(data string) (RequestLine, error) {
	// Split on first newline to get request line
	parts := strings.SplitN(data, "\r\n", 2)
	if len(parts) == 0 {
		return RequestLine{}, errors.New("empty request")
	}

	// Split request line into components
	requestParts := strings.Split(parts[0], " ")
	if len(requestParts) != 3 {
		return RequestLine{}, errors.New("invalid request line format")
	}

	method := requestParts[0]
	target := requestParts[1]
	version := requestParts[2]

	// Validate method
	if !isValidMethod(method) {
		return RequestLine{}, errors.New("invalid method")
	}

	// Validate version
	if version != "HTTP/1.1" {
		return RequestLine{}, errors.New("unsupported HTTP version")
	}

	return RequestLine{
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
