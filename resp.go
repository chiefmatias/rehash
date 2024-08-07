package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type RespMessage struct {
	typ     byte
	integer int
	str     string
	values  []RespMessage
}

const (
	BlobString   = byte('$')
	SimpleString = byte('+')
	SimpleErr    = byte('-')
	Integer      = byte(':')
	Array        = byte('*')
)

type reader func(i *bufio.Reader) (RespMessage, error)

var readers = [256]reader{}

func init() {
	readers[BlobString] = handleBlobString
	readers[SimpleString] = handleSimpleString
	readers[Integer] = handleInteger
	readers[Array] = handleArray
}

func handleBlobString(reader *bufio.Reader) (msg RespMessage, err error) {
	msg.typ = BlobString
	msg.integer, err = readInteger(reader)
	if err != nil {
		return msg, err
	}
	msg.str, err = readBulk(reader, msg.integer)
	return msg, err
}

func handleArray(reader *bufio.Reader) (msg RespMessage, err error) {
	msg.typ = Array
	msg.integer, err = readInteger(reader)
	if err != nil {
		return msg, err
	}
	msg.values = make([]RespMessage, 0, msg.integer)
	_, err = readArray(reader, &msg)

	return msg, err
}
func handleSimpleString(reader *bufio.Reader) (msg RespMessage, err error) {
	msg.typ = SimpleString
	msg.str, err = readSimple(reader)

	return msg, err
}

func handleInteger(reader *bufio.Reader) (msg RespMessage, err error) {
	msg.typ = Integer
	msg.integer, err = readInteger(reader)

	return msg, err
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
	return integer, nil
}

func readSimple(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	str := strings.TrimSuffix(string(data), "\r\n")

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

	return str, nil
}

func readArray(reader *bufio.Reader, msg *RespMessage) (*RespMessage, error) {
	for i := 0; i < msg.integer; i++ {
		item, err := respParser(reader)
		if err != nil {
			return msg, err
		}
		msg.values = append(msg.values, item)
	}
	return msg, nil
}

func respParser(reader *bufio.Reader) (RespMessage, error) {
	typ, err := reader.ReadByte()

	if err != nil {
		return RespMessage{}, err
	}

	handler := readers[typ]
	if handler == nil {
		return RespMessage{}, fmt.Errorf("unsupported message type: %c", typ)
	}

	return handler(reader)
}
