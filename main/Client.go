package main

import (
	"net"
	"encoding/gob"
)

var nodeConnMap = make(map[string]net.Conn)

func connect(name string) {
	_, ok := nodeConnMap[name]
	if !ok {
		conn, _ := net.Dial("tcp", NodeMap[name].address+":" + NodeMap[name].port)
		nodeConnMap[name] = conn
	}
}

func SendSocket(message TcpMessage) {
	dest := message.Destination
	connect(dest)
	enc := gob.NewEncoder(nodeConnMap[dest])
	enc.Encode(message)
	nodeConnMap[dest].Close()
}