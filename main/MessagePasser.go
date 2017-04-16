package main

import (
	"strings"
	"net/url"
	"reflect"
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

func Multicast(origSender string, groupName string, msgPurposeUrl string, reqParam url.Values, msgData string) {
	// Multicast to all of the drones and let them determine whether to deliver or not
	keys := reflect.ValueOf(Swarm).MapKeys()
	drones := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		drones[i] = keys[i].String()
	}
	for _, droneId := range drones {
		url := "http://" + DroneNodeMap[droneId] + msgPurposeUrl + "?type=" + MULTICAST_TYPE + "&" + reqParam.Encode()
		msg := MulticastMessage{origSender, groupName, seqNum, msgData}
		makeGetRequest(url, toJsonString(msg))
	}
	if strings.Compare(origSender, DroneId) == 0 {
		seqNum += 1
	}
}

