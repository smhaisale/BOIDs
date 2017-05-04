package main

import (
	"net/url"
	"strconv"
)

/*
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
*/

// origSender use this function
func multicastMaekawa(origSender string, dest string, purposeUrl string, msg MaekawaMessage, seqNum int) {
	// Multicast to all of the drones and let them determine whether to deliver or not
	reqParam := url.Values{};
	reqParam.Add("origSender", origSender)
	reqParam.Add("dest", dest)
	reqParam.Add("seqNum", strconv.Itoa(seqNum))
	asyncGetRequest("http://" + drone.Address + purposeUrl + "?" + reqParam.Encode(), toJsonString(msg))
	//makeGetRequest("http://" + drone.Address + purposeUrl + "?" + reqParam.Encode(), toJsonString(msg))
	for _, otherDrone := range swarm {
		url := "http://" + otherDrone.Address + purposeUrl + "?" + reqParam.Encode()
		//makeGetRequest(url, toJsonString(msg))
		asyncGetRequest(url, toJsonString(msg))
	}
}

// not origSender use this function
func sendMulticast(purposeUrl string, reqParam url.Values, msg MaekawaMessage) {
	makeGetRequest("http://" + drone.Address + purposeUrl + "?"+ reqParam.Encode(), toJsonString(msg))
	for _, otherDrone := range swarm {
		url := "http://" + otherDrone.Address + purposeUrl + "?"+ reqParam.Encode()
		//makeGetRequest(url, toJsonString(msg))
		asyncGetRequest(url, toJsonString(msg))
	}
}
