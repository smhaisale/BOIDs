package main

import "fmt"

func Send(message TcpMessage) {
	SendSocket(message)
}

func Receive() {
	message := ReceiveSocket()
	fmt.Println(message)
}
