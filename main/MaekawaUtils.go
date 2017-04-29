package main

import (
	"log"
	"fmt"
	"strconv"
	"reflect"
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
//var seqNum map[string]int = make(map[string]int)


func getDrones() []string{
	keys := reflect.ValueOf(swarm).MapKeys()
	drones := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		drones[i] = keys[i].String()
	}
	return drones
}

// todo
func formPermGroup() {
	permissionGroup = getDrones()
}

func request(path PathLock) {
	myPathLock = path
	seqNum += 1
	formPermGroup()
	for _, otherDroneId := range permissionGroup {
		/*
		_, exist := seqNum[otherDroneId]
		if !exist {
			seqNum[otherDroneId] = 0
		}
		seqNum[otherDroneId] += 1
		*/
		log.Println("Request seqNum " + strconv.Itoa(seqNum))
		reqMsg := MaekawaMessage{drone.ID, otherDroneId, REQUEST, path}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, reqMsg, seqNum)
	}
}

func release() {
	//todo
	ackNo = 0
	seqNum += 1
	for _, otherDroneId := range permissionGroup {
		/*
		_, exist := seqNum[otherDroneId]
		if !exist {
			seqNum[otherDroneId] = 0
		}
		seqNum[otherDroneId] += 1
		*/
		log.Println("Release seqNum " + strconv.Itoa(seqNum))
		rlsMsg := MaekawaMessage{drone.ID, otherDroneId, RELEASE, myPathLock}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, rlsMsg, seqNum)
	}
}

func ack(dest string) {
	/*
	_, exist := seqNum[dest]
	if !exist {
		seqNum[dest] = 0
	}
	seqNum[dest] += 1
	*/
	seqNum += 1
	log.Println("Ack seqNum " + strconv.Itoa(seqNum))
	ackMsg := MaekawaMessage{drone.ID, dest, ACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, ackMsg, seqNum)
}

func nack(dest string) {
	/*
	_, exist := seqNum[dest]
	if !exist {
		seqNum[dest] = 0
	}
	seqNum[dest] += 1
	*/
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
