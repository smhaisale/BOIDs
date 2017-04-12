package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

var recvQueue = make(chan TcpMessage)

func handleConnection(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		fmt.Println("accept")
		if err != nil {
		}
		decoder := gob.NewDecoder(conn)
		msg := &TcpMessage{}
		decoder.Decode(msg)
		fmt.Println(*msg)
		if msg != nil {
			recvQueue <- *msg
		}
		conn.Close()
	}
}
func ReceiveSocket() TcpMessage {
	return <-recvQueue
}

func Listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	fmt.Println("listen")
	if err != nil {
		return
	}
	go handleConnection(ln)
	fmt.Println("done")
}
