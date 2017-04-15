package main

import "fmt"

func main() {
    var name string
    fmt.Scanf("%s", &name)

    Listen(NodeMap[name].port)
    go Receive()

    myDrone := DroneObject{Position{0, 1, 2}, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, Speed{1, 2, 3}}
    msgData := MessageData{[]Drone{Drone{"drone1", "", "", myDrone}}}
    var dest string
    for {
        fmt.Println("Input Destination:")
        fmt.Scanf("%s", &dest)
        fmt.Println("Sending message to " + dest)
        msg := TcpMessage{name, dest, 0, true, "Message", VectorTime{map[string]int{"alice": 1}}, msgData, MulticastData{}}
        Send(msg)
    }
}