package headers

import (
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func isValidHeaderKeyChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) ||
		strings.ContainsRune("!#$%&'*+-.^_`|~", r)
}

func validateHeaderKey(key string) error {
	for _, r := range key {
		if !isValidHeaderKeyChar(r) {
			return errors.New("invalid character in header key")
		}
	}
	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Check for empty data
	if len(data) == 0 {
		return 0, false, nil
	}

	// Check for end of headers (CRLF at start)
	if len(data) >= 2 && data[0] == '\r' && data[1] == '\n' {
		return 2, true, nil
	}

	// Find the end of the header line
	end := strings.Index(string(data), "\r\n")
	if end == -1 {
		return 0, false, nil
	}

	// Split the header line
	headerLine := string(data[:end])
	colonIndex := strings.Index(headerLine, ":")
	if colonIndex == -1 {
		return 0, false, errors.New("invalid header format: missing colon")
	}

	key := headerLine[:colonIndex]
	if len(key) == 0 || key[len(key)-1] == ' ' {
		return 0, false, errors.New("invalid header format: empty key or trailing space")
	}

	// Extract and clean key and value
	key = strings.TrimSpace(headerLine[:colonIndex])
	value := strings.TrimSpace(headerLine[colonIndex+1:])

	// Validate key format
	if err := validateHeaderKey(key); err != nil {
		return 0, false, err
	}

	// Add to headers (convert key to lowercase)
	h[strings.ToLower(key)] = value

	// Return bytes consumed (header line + CRLF)
	return end + 2, false, nil
}
