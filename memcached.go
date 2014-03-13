package main

import (
	"net"
)

func handleConnection(conn net.Conn) {
	println("Yahoo!")
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
