package main

import (
	"fmt"
	"net"
	"os"
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
			args := argsParser(receivedData)
			command := strings.ToUpper(args[0])

			fmt.Printf("Command: '%s'\n", command)

			fmt.Println("Args:", args)
			fmt.Println("len(Args):", len(args))

			if command == "PING" {
				_, err = conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			}
			if command == "QUIT" {
				fmt.Println("Quitting connection.")
				os.Exit(1)
			}

			receivedData.Reset()

		}
	}
}

func argsParser(receivedData strings.Builder) []string {
	args := strings.Split(receivedData.String(), " ")

	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}

	return args

}
