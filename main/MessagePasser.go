package main

import "fmt"

func Send(message TcpMessage) {
	SendSocket(message)
}

func Receive() {
	for {
		recvMsg := ReceiveSocket()
		fmt.Println("Receive Message from " + recvMsg.Source)
		fmt.Println(recvMsg)
	}
}
