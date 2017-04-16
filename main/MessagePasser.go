package main

import (
	"net/url"
	"reflect"
	"strings"
)

/**
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
**/

var seqNum int = 0

func Multicast(origSender string, dest string, msgPurposeUrl string, reqParam url.Values, msgData string) {
	// Multicast to all of the drones and let them determine whether to deliver or not
	keys := reflect.ValueOf(Swarm).MapKeys()
	drones := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		drones[i] = keys[i].String()
	}
	for _, droneId := range drones {
		url := "http://" + DroneNodeMap[droneId] + msgPurposeUrl + reqParam.Encode()
		msg := MulticastMessage{origSender, dest, seqNum, msgData}
		makeGetRequest(url, toJsonString(msg))
	}
	if strings.Compare(origSender, drone.ID) == 0 {
		seqNum += 1
	}
}
