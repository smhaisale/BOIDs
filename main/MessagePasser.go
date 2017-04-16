package main

import (
	"fmt"
	"strings"
	"net/url"
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
	for _, droneId := range GroupMap[ALLDRONES_GROUP] {
		url := "http://" + DroneNodeMap[droneId] + msgPurposeUrl + "?type=" + MULTICAST_TYPE + "&" + reqParam.Encode()
		msg := MulticastMessage{origSender, groupName, seqNum, msgData}
		makeGetRequest(url, toJsonString(msg))
	}
	if strings.Compare(origSender, DroneId) == 0 {
		seqNum += 1
	}
}

