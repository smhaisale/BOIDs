package main

import (
	"fmt"
	"net"
	"bufio"
)

var recvQueue = make(chan TcpMessage)

func handleConnection(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		rw := bufio.NewReader(conn)
		msgString, _ := rw.ReadString('\n')
		var msg TcpMessage
		fromJsonString(&msg, string(msgString))
		recvQueue <- msg
		conn.Close()
	}
}

func ReceiveSocket() TcpMessage {
	msg := <-recvQueue
	return msg
}

func Listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Launch the server...")
	go handleConnection(ln)
}
