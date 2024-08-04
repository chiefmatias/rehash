package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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
	defer conn.Close()
	reader := bufio.NewReader(conn)
	respParser(reader)
}

func respParser(reader *bufio.Reader) {
	char, _ := reader.ReadByte()

	if char != '$' {
		fmt.Println("This is not a bulk string!")
	}

	sizeBytes, _ := reader.ReadBytes('\n')
	size, _ := strconv.Atoi(strings.TrimSpace(string(sizeBytes)))

	fmt.Println("size is:", size)

	data := make([]byte, size)
	_, err := reader.Read(data)
	if err != nil {
		fmt.Println("Something happened with the reader.", err)
	}

	fmt.Println("Data is:", string(data))
}
