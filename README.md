# MyHttpProtocol

A Go implementation of an HTTP/1.1 protocol parser with streaming support.

## Features

- HTTP/1.1 request line parsing
- Streaming data support
- Memory-efficient buffer management
- Stateful parsing
- Support for standard HTTP methods

## Structure

```
.
├── internal/
│   └── request/
│       ├── request.go      # HTTP request parsing implementation
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

## License

MIT