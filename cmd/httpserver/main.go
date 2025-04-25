package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"chillhttp/internal/request"
	"chillhttp/internal/response"
	"chillhttp/internal/server"

	"github.com/pingcap/log"
)

const port = 42069

func HttpHandler(w *response.Writer, req *request.Request) {
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