package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Error creating UDP connection: %v\n", err)
		return
	}

	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)
	if reader == nil {
		fmt.Println("Error creating reader")
		return
	}

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading line: %v\n", err)
			return
		}

		udpConn.Write([]byte(line))
	}
}
