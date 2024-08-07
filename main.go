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

	msg, err := respParser(reader)
	if err != nil {
		fmt.Println("An awful error just occurred!", err)
		return
	}
	fmt.Println("This is the message:", msg)

}

type RespMessage struct {
	typ     byte
	integer int
	str     string
	values  []RespMessage
}

const (
	typeBlobString   = byte('$')
	typeSimpleString = byte('+')
	typeSimpleErr    = byte('-')
	typeInteger      = byte(':')
	typeArray        = byte('*')
)

func respParser(reader *bufio.Reader) (RespMessage, error) {
	var err error

	msg := &RespMessage{}
	msg.typ, err = reader.ReadByte()
	if err != nil {
		return *msg, err
	}

	fmt.Printf("Reading type: %c\n", msg.typ)

	switch msg.typ {
	case typeInteger:
		msg.integer, err = readInteger(reader)

	case typeSimpleString:
		msg.str, err = readSimple(reader)

	case typeBlobString:
		msg.integer, err = readInteger(reader)
		if err != nil {
			return *msg, err
		}
		msg.str, err = readBulk(reader, msg.integer)

	case typeArray:
		msg.integer, err = readInteger(reader)
		if err != nil {
			return *msg, err
		}
		msg.values = make([]RespMessage, 0, msg.integer)
		_, err = readArray(reader, msg)

	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Parsed message: %+v\n", *msg)
	return *msg, err
}

func readInteger(reader *bufio.Reader) (int, error) {
	sizeBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}

	integer, err := strconv.Atoi(strings.TrimSpace(string(sizeBytes)))
	if err != nil {
		return 0, err
	}

	fmt.Println("size is:", integer)
	return integer, nil
}

func readSimple(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}

	str := strings.TrimSuffix(string(data), "\r\n")
	fmt.Println("string is:", string(data))

	return str, nil
}

func readBulk(reader *bufio.Reader, size int) (string, error) {
	data := make([]byte, size)
	_, err := reader.Read(data)
	if err != nil {
		return "", nil
	}

	str := strings.TrimSuffix(string(data), "\r\n")

	if _, err := reader.ReadBytes('\n'); err != nil {
		return "", err
	}

	fmt.Println("bulk string is:", string(str))
	return str, nil
}

func readArray(reader *bufio.Reader, msg *RespMessage) (*RespMessage, error) {
	for i := 0; i < msg.integer; i++ {
		fmt.Printf("Reading array element %d/%d\n", i+1, msg.integer)
		item, err := respParser(reader)
		if err != nil {
			return msg, err
		}
		msg.values = append(msg.values, item)
		fmt.Printf("Appended item: %+v\n", item)
	}
	return msg, nil
}
