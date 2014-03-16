package main

import (
	"github.com/takuji/memcached/protocol"
	"log"
	"net"
	"strings"
)

type Command struct {
	key string
}

func newCommand(s string) (*Command, error) {
	println(s)
	elms := strings.Split(s, " ")
	command := new(Command)
	command.key = elms[0]
	return command, nil
}

func (c *Command) handle(s string) {
	println(s)
}

func handleConnection(conn net.Conn) {
	r := protocol.NewRequestReader(conn)

	for {
		req, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", req)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":11211")
	if err != nil {
		println(err)
	} else {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn)
		}
	}
}
