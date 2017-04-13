package main

import (
	"fmt"
	"net"
	"encoding/gob"
)

var recvQueue = make(chan TcpMessage)

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	var msg TcpMessage
	dec.Decode(&msg)
	recvQueue <- msg
}

func ReceiveSocket() TcpMessage {
	msg := <-recvQueue
	fmt.Println(msg)
	return msg
}

func Listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		fmt.Println("accept")
		if err != nil {
		}
		go handleConnection(conn)
	}
}
