package main

import "fmt"

func main() {
	var name string
	fmt.Scanf("%s", &name)
	go Listen(NodeMap[name].port)
	go Receive()
	fmt.Println("start")
	myDrone := Drone{"0", Position{0, 1, 2}, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, Speed{1, 2, 3}}
	msgData := MessageData{[]Drone{myDrone}}
	var dest string
	for {
		fmt.Scanf("%s", &dest)
		fmt.Println("send message")
		msg := TcpMessage{name, dest, 0, true, "Message", VectorTime{map[string]int{"alice": 1}}, msgData, MulticastData{}}
		Send(msg)
	}

}
