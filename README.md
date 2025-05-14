# Chill Http

A Go implementation of an HTTP/1.1 server with streaming support, chunked encoding, and proxy capabilities.

## Features

- HTTP/1.1 request line parsing
- Streaming data support
- Memory-efficient buffer management
- Stateful parsing
- Support for standard HTTP methods
- Custom status code handling with defined constants
- Chunked transfer encoding support
- HTTP proxy functionality
- Response trailers support
- Custom response writer implementation

## Structure

```
.
├── internal/
│   └── request/
│       ├── request.go      # HTTP requests parsing implementation
│       └── request_test.go # Test cases for request parsing
└── cmd/
    └── udpsender/
    |   └── main.go         # UDP client for testing
    └── tcplistener/
    |   └── main.go         # TCP client for testing
    └── httpserver/
        └── main.go         # hhtp server for testing
```

## Usage

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

## Implementation Details

The parser implements:
- Streaming data handling
- Buffer management with growth
- State tracking (initialized/done)
- HTTP/1.1 request line validation
- Method validation (uppercase letters only)
- HTTP response status code handling with defined constants (200, 400, 500)

## Examples

### Basic Requests

```bash
# Success Response
$ curl http://localhost:42069/
<html>
    <head>
        <title>200 OK</title>
    </head>
    <body>
        <h1>Success!</h1>
        <p>Your request was an absolute banger.</p>
    </body>
</html>

# Bad Request
$ curl http://localhost:42069/yourproblem
<html>
    <head>
        <title>400 Bad Request</title>
    </head>
    <body>
        <h1>Bad Request</h1>
        <p>Your request honestly kinda sucked.</p>
    </body>
</html>

# Server Error
$ curl http://localhost:42069/myproblem
<html>
    <head>
        <title>500 Internal Server Error</title>
    </head>
    <body>
        <h1>Internal Server Error</h1>
        <p>Okay, you know what? This one is on me.</p>
    </body>
</html>
```

### Proxy Requests with Chunked Transfer

```bash
# View raw chunked response (using netcat)
$ nc localhost 42069
GET /httpbin/stream/100 HTTP/1.1
Host: localhost

HTTP/1.1 200 OK
Transfer-Encoding: chunked
Trailer: X-Content-Sha256, X-Content-Length
[...]

# Each chunk is preceded by its length in hexadecimal
a4
{"args": {}, "headers": {...}, "origin": "...", "url": "https://httpbin.org/stream/100"}
a4
{"args": {}, "headers": {...}, "origin": "...", "url": "https://httpbin.org/stream/100"}
...
0

X-Content-SHA256: [hash of the complete response]
X-Content-Length: [total bytes transferred]
```

### Video Streaming

```bash
# Stream video content
$ curl http://localhost:42069/video -o video.mp4
```

### Using Different HTTP Methods

The server supports standard HTTP methods (GET, POST, etc.) with proper validation:
```bash
# GET request
$ curl http://localhost:42069/

# POST request
$ curl -X POST http://localhost:42069/

# Invalid method
$ curl -X INVALID http://localhost:42069/
HTTP/1.1 400 Bad Request
```

## License

MIT
