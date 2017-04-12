package main

import (
	"encoding/gob"
	"net"
	"fmt"
)

var nodeConnMap = make(map[string]net.Conn)

func connect(name string) {
	_, ok := nodeConnMap[name]
	if !ok {
		conn, _ := net.Dial("tcp", NodeMap[name].address+":"+NodeMap[name].port)
		nodeConnMap[name] = conn
	}
}

func SendSocket(message TcpMessage) {
	dest := message.destination
	connect(dest)
	fmt.Println(nodeConnMap[dest])
	encoder := gob.NewEncoder(nodeConnMap[dest])
	encoder.Encode(message)
}