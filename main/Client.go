package main

import (
	"net"
	"fmt"
)

var nodeConnMap = make(map[string]net.Conn)

func connect(name string) {
	_, ok := nodeConnMap[name]
	if !ok {
		conn, err := net.Dial("tcp", NodeMap[name].address+":"+NodeMap[name].port)
		if err != nil {
			fmt.Println(err)
		}
		nodeConnMap[name] = conn
	}
}

func SendSocket(message TcpMessage) {
	dest := message.Destination
	connect(dest)
	fmt.Fprint(nodeConnMap[dest], toJsonString(message)+"\n")
}