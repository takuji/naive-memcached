package main

import (
	"github.com/takuji/memcached/protocol"
	"log"
	"net"
	"strings"
	"time"
)

var (
	data map[string]*Item = make(map[string]*Item)
)

func init() {
	log.Printf("%+v", data)
}

type Item struct {
	Value     string
	CreatedAt time.Time
}

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
			log.Println(err)
			continue
		}
		log.Printf("%+v", req)
		switch r := req.(type) {
		case *protocol.SetRequest:
			log.Println("SetRequest!")
			log.Println(r.Key)
			data[r.Key] = &Item{Value: r.Data, CreatedAt: time.Now()}
			log.Printf("Current items: %+v", len(data))
		default:
			log.Println("???")
		}
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
