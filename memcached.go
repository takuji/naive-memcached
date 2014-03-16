package main

import (
	"github.com/takuji/memcached/protocol"
	"log"
	"net"
	"time"
)

var (
	data map[string]*Item = make(map[string]*Item)
)

type Item struct {
	Value     string
	CreatedAt time.Time
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
		case *protocol.GetRequest:
			log.Println("GetRequest!")
			log.Println(r.Key)
			item := data[r.Key]
			conn.Write([]byte(item.Value))
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
