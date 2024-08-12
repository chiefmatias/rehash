package main

import (
	"fmt"
	"strings"
)

const (
	Ping = "PING"
	Echo = "ECHO"
)

type command func(msg RespMessage) (RespMessage, error)

var commands = make(map[string]command)

func init() {
	commands[Ping] = pingCommand
	commands[Echo] = echoCommand
}

func commandHandler(msg RespMessage) (RespMessage, error) {
	var cmd string

	if msg.typ == Array {
		cmd = msg.values[0].str
	} else {
		cmd = msg.str
	}

	cmd = strings.ToUpper(cmd)
	handler := commands[cmd]
	if handler == nil {
		str := fmt.Sprintf("ERR unknown command '%s'", cmd)
		return RespMessage{typ: SimpleErr, str: str}, fmt.Errorf("unsupported command: %s", cmd)
	}

	return handler(msg)
}

func pingCommand(msg RespMessage) (answer RespMessage, err error) {
	answer.typ = SimpleString
	answer.str = "PONG"

	return answer, err
}

func echoCommand(msg RespMessage) (answer RespMessage, err error) {
	if msg.integer == 2 {
		answer.typ = BulkString
		answer.str = msg.values[1].str
		answer.integer = msg.values[1].integer
	} else {
		answer.typ = SimpleErr
		answer.str = fmt.Sprintf("ERR invalid ammount of arguments:'%d'", msg.integer)
		err = fmt.Errorf("invalid ammount of arguments: %d", msg.integer)
	}

	return answer, err
}
