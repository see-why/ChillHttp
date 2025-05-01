package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"chillhttp/internal/request"
	"chillhttp/internal/response"
	"chillhttp/internal/server"

	"github.com/pingcap/log"
)

const port = 42069

func HttpHandler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
		targetURL := "https://httpbin.org" + proxyPath

		resp, err := http.Get(targetURL)
		if err != nil {
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(response.GetDefaultHeaders(0))
			w.WriteBody([]byte("Proxy error\n"))
			return
		}
		defer resp.Body.Close()

		// Remove Content-Length, add Transfer-Encoding: chunked
		headers := response.GetDefaultHeaders(0)
		delete(headers, "Content-Length")
		headers["Transfer-Encoding"] = "chunked"
		headers["Trailer"] = "X-Content-Sha256, X-Content-Length"

		// Write status line and headers
		w.WriteStatusLine(response.StatusCode(resp.StatusCode))
		w.WriteHeaders(headers)
		
		var fullBody []byte
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				// Print n for debugging
				chunk := buf[:n]
				fmt.Println("Read bytes:", n)
				fullBody = append(fullBody, chunk...)
				w.WriteChunkedBody(chunk)
			}
			if err != nil {
				break
			}
		}
		w.WriteChunkedBodyDone()

		hashSum := sha256.Sum256(fullBody)
		trailers := response.GetDefaultTrailerHeaders(len(fullBody), fmt.Sprintf("%x", hashSum))
		w.WriteTrailers(trailers)

		return
	}

	var body []byte
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.BadRequest)
		body = []byte(`<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
		</html>`)
	case "/myproblem":
		w.WriteStatusLine(response.InternalServerError)
		body = []byte(`<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
		</html>`)
	default:
		w.WriteStatusLine(response.OK)
		body = []byte(`<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
		</html>`)
	}

	header := response.GetDefaultHeaders(len(body))
	header["Content-Type"] = "text/html"
	err := w.WriteHeaders(header)
	if err != nil {
		fmt.Println("Error writing headers: ", err)
		return
	}

	w.Writer.Write(body)
}

func main() {
	server, err := server.Serve(port, HttpHandler)
	if err != nil {
		fmt.Println("Error starting server: %w", err)
	}
	defer server.Close()
	fmt.Println("Server listening on :", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Info("Received shutdown signal, shutting down server...")
}