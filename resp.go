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

func respParser(reader *bufio.Reader) (RespMessage, error) {
	var err error

	msg := &RespMessage{}
	msg.typ, err = reader.ReadByte()
	if err != nil {
		return *msg, err
	}

	switch msg.typ {
	case Integer:
		msg.integer, err = readInteger(reader)

	case SimpleString:
		msg.str, err = readSimple(reader)

	case BlobString:
		msg.integer, err = readInteger(reader)
		if err != nil {
			return *msg, err
		}
		msg.str, err = readBulk(reader, msg.integer)

	case Array:
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
	fmt.Println("Parsed message:", msg)
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
