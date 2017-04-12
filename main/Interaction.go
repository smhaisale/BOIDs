package main

import "fmt"

func main() {
	var name string
	fmt.Scanf("%s", &name)
	Listen(NodeMap[name].port)
	fmt.Println("listening")
	myDrone := Drone{Position{0, 1, 2}, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, Speed{1, 2, 3}}
	msgData := MessageData{[]Drone{myDrone}}
	msg := TcpMessage{"alice", "bob", 0, false, "Message", VectorTime{map[string]int{"alice": 1}}, msgData, MulticastData{}}
	var input int
	fmt.Scanf("%d", &input)
	if input == 0 {
		fmt.Println("send message")
		Send(msg)
	}
	go Receive()

}
