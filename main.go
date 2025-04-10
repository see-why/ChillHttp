package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)
	currentLine := ""
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		// Split the current chunk into parts based on newlines
		parts := strings.Split(string(buffer[:n]), "\n")

		// Process all parts except the last one
		for i := 0; i < len(parts)-1; i++ {
			currentLine += parts[i]
			fmt.Printf("read: %s\n", currentLine)
			currentLine = ""
		}

		// Add the last part to our current line
		currentLine += parts[len(parts)-1]
	}

	// Print any remaining content in currentLine
	if currentLine != "" {
		fmt.Printf("read: %s\n", currentLine)
	}
}
