package main

import (
	"bufio"
	"fmt"
	"net"
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

	for {

		msg, err := respParser(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Parsed message:", msg)

		answer, _ := respSerializer(RespMessage{typ: SimpleString, str: "OK"})
		conn.Write(answer)

	}
}
