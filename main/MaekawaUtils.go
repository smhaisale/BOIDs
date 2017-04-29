package main

import (
	"log"
	"fmt"
	"strconv"
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

var seqNum int = 0

// todo
func formPermGroup() {
	for k, _ := range swarm {
		permissionGroup = append(permissionGroup, k)
	}
}

func request(path PathLock) {
	myPathLock = path
	//todo
	formPermGroup()
	for _, otherDroneId := range permissionGroup {
		seqNum += 1
		log.Println("Request seqNum " + strconv.Itoa(seqNum))
		reqMsg := MaekawaMessage{drone.ID, otherDroneId, REQUEST, path}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, reqMsg, seqNum)
	}
}

func release() {
	//todo
	ackNo = 0
	formPermGroup()
	for _, otherDroneId := range permissionGroup {
		seqNum += 1
		log.Println("Release seqNum " + strconv.Itoa(seqNum))
		rlsMsg := MaekawaMessage{drone.ID, otherDroneId, RELEASE, myPathLock}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, rlsMsg, seqNum)
	}
}

func ack(dest string) {
	seqNum += 1
	log.Println("Ack seqNum " + strconv.Itoa(seqNum))
	ackMsg := MaekawaMessage{drone.ID, dest, ACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, ackMsg, seqNum)
}

func nack(dest string) {
	seqNum += 1
	log.Println("Nack seqNum " + strconv.Itoa(seqNum))
	nackMsg := MaekawaMessage{drone.ID, dest, NACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, nackMsg, seqNum)
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
	log.Println("---------")
	log.Println("path request queue: ")
	for i,_ := range pathRequestQueue {
		log.Println(i)
	}
	log.Println("---------")
	log.Println("---------")
	log.Println("curr path lock list: ")
	for i,_ := range currPathLockList {
		log.Println(i)
	}
	log.Println("---------")
	if hasIntersect {
		log.Println("has intersect")
		pathRequestQueue[source] = path
		nack(source)
	} else {
		log.Println("no intersect")
		currPathLockList[source] = path
		ack(source)
	}
}

func handleRelease(msg MaekawaMessage) {
	source := msg.Source
	_, exist := currPathLockList[source]
	if exist {
		log.Println("delete currPathlock")
		delete(currPathLockList, source)
	}
	for reqSource, reqPathLock := range pathRequestQueue {
		hasIntersect := false
		for _, currPathLock := range currPathLockList {
			if isIntersect(reqPathLock, currPathLock) {
				hasIntersect = true
				break
			}
		}
		if !hasIntersect {
			currPathLockList[reqSource] = reqPathLock
			delete(pathRequestQueue, reqSource)
			ack(reqSource)
			break;
		}
	}
}

func handleAck(msg MaekawaMessage) {
	ackNo += 1
	if ackNo >= len(permissionGroup) {
		//todo enter cs
		fmt.Println("Enter CS");
		var stop string
		fmt.Scanln(&stop);
		release();
	}
}

func handleNack(msg MaekawaMessage) {

}

func isIntersect(path1 PathLock, path2 PathLock) bool {
	return (dist3D_Segment_to_Segment(path1, path2) < 2)
}
