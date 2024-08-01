package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Waiting for connection...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Accepted connection from:", conn.RemoteAddr())

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	var receivedData strings.Builder
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		receivedData.WriteString(string(buffer[:n]))
		if strings.Contains(receivedData.String(), "\r\n") {
			command := strings.TrimSpace(receivedData.String())
			fmt.Printf("Command: '%s'\n", command)

			if command == "PING" {
				_, err = conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			}
			if command == "QUIT" {
				fmt.Println("Quitting connection.")
				conn.Close()
			}

			receivedData.Reset()

		}
	}
}
