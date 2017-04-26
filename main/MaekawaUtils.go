package main

import (
	"log"
	"net/http"
)

type PathLock struct {
	From Position
	To   Position
}

var permissionGroup []string = make([]string, len(swarm))
var currPathLockList map[string]PathLock = make(map[string]PathLock)
var pathRequestQueue map[string]PathLock = make(map[string]PathLock)
var myPathLock PathLock
var ackNo int

type MaekawaMessage struct {
	Source      string
	Destination string
	Type        string
	Path        PathLock
}

var REQUEST = "REQUEST"
var RELEASE = "RELEASE"
var ACK = "ACK"
var NACK = "NACK"

// todo
func formPermGroup() {
	for k, _ := range swarm {
		permissionGroup = append(permissionGroup, k)
	}
}

func request(path PathLock) {
	myPathLock = path
	for _, otherDroneId := range permissionGroup {
		requestUrl := "http://" + swarm[otherDroneId].Address + DRONE_MAEKAWA_MESSAGE_URL
		reqMsg := MaekawaMessage{drone.ID, otherDroneId, REQUEST, path}
		makeGetRequest(requestUrl, toJsonString(reqMsg))
	}
}

func release() {
	for _, otherDroneId := range permissionGroup {
		releaseUrl := "http://" + swarm[otherDroneId].Address + DRONE_MAEKAWA_MESSAGE_URL
		rlsMsg := MaekawaMessage{drone.ID, otherDroneId, RELEASE, myPathLock}
		makeGetRequest(releaseUrl, toJsonString(rlsMsg))
	}
}

func ack(dest string) {
	ackUrl := "http://" + swarm[dest].Address + DRONE_MAEKAWA_MESSAGE_URL
	ackMsg := MaekawaMessage{drone.ID, dest, ACK, PathLock{}}
	makeGetRequest(ackUrl, toJsonString(ackMsg))
}

func nack(dest string) {
	nackUrl := "http://" + swarm[dest].Address + DRONE_MAEKAWA_MESSAGE_URL
	nackMsg := MaekawaMessage{drone.ID, dest, NACK, PathLock{}}
	makeGetRequest(nackUrl, toJsonString(nackMsg))
}

func handleRequest(msg MaekawaMessage) {
	source := msg.Source
	path := msg.Path
	hasIntersect := false
	for _, pathLock := range currPathLockList {
		if isIntersect(path, pathLock) {
			hasIntersect = true
			break
		}
	}
	if hasIntersect {
		pathRequestQueue[source] = path
		nack(source)
	}
}

func handleRelease(msg MaekawaMessage) {
	ackNo = 0
	source := msg.Source
	delete(currPathLockList, source)
	for reqSource, reqPathLock := range pathRequestQueue {
		hasIntersect := false
		for _, currPathLock := range currPathLockList {
			if isIntersect(reqPathLock, currPathLock) {
				hasIntersect = true
				break
			}
		}
		if !hasIntersect {
			ack(source)
			currPathLockList[reqSource] = reqPathLock
			delete(pathRequestQueue, reqSource)
		}
	}
}

func handleAck(msg MaekawaMessage) {
	ackNo += 1
	if ackNo >= len(permissionGroup) {
		//todo enter cs
	}
}

func handleNack(msg MaekawaMessage) {

}

func isIntersect(path1 PathLock, path2 PathLock) bool {
	return (dist3D_Segment_to_Segment(path1, path2) < 2)
}

/*
func createLockMessage(data string) (message MaekawaMessage) {
    reset()
    message.Source = drone.ID
    message.Data = data

    log.Println("Created prepare message: ", message)
    return message
}

func createReleaseMessage(dest string, data string) (message MaekawaMessage) {
    message.Destination = dest
    message.Data = data

    log.Println("Created promise message: ", message)
    return message
}

func createRejectMessage(data string) (message MaekawaMessage) {
    message.Source = drone.ID
    message.Data = data

    log.Println("Created accept message: ", message)
    return message
}

func createFailedMessage(message MaekawaMessage) {

}

func handleMaekawaMessage(w http.ResponseWriter, r *http.Request) {
    message := MaekawaMessage{}
    getRequestBody(&message, r)
}
*/
