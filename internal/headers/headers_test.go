package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestValidSingleHeaderWithCapitalLetters(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-Type: application/json\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 32, n)
	assert.False(t, done)
}

func TestValidSingleHeaderWithExtraWhitespace(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:    localhost:42069    \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 30, n)
	assert.False(t, done)
}

func TestValidTwoHeadersWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	headers["existing"] = "value"

	data := []byte("Host: localhost:42069\r\nContent-Type: application/json\r\n\r\n")

	// Parse first header
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Parse second header
	n, done, err = headers.Parse(data[23:])
	require.NoError(t, err)
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 32, n)
	assert.False(t, done)
}

func TestValidDone(t *testing.T) {
	headers := NewHeaders()
	data := []byte("\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestInvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestInvalidCharacterInHeaderKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestValidSpecialCharactersInHeaderKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Forwarded-For: 127.0.0.1\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", headers["x-forwarded-for"])
	assert.Equal(t, 28, n)
	assert.False(t, done)
}
