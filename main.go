package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(reader io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer reader.Close()

		buffer := make([]byte, 8)
		currentLine := ""

		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Error reading file: %v\n", err)
				return
			}

			parts := strings.Split(string(buffer[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				ch <- currentLine
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]
		}

		if currentLine != "" {
			ch <- currentLine
		}
	}()

	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}
