# ChillHttp

A Go implementation of an HTTP/1.1 protocol parser with streaming support.

## Features

- HTTP/1.1 request line parsing
- Streaming data support
- Memory-efficient buffer management
- Stateful parsing
- Support for standard HTTP methods
- Custom status code handling with defined constants

## Structure

```
.
├── internal/
│   └── request/
│       ├── request.go      # HTTP requests parsing implementation
│       └── request_test.go # Test cases for request parsing
└── cmd/
    └── udpsender/
        └── main.go         # UDP client for testing
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

## License

MIT