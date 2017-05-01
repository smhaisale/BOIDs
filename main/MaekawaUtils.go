package main

import (
	"log"
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

var seqNum = map[string]int{REQUEST: 0, RELEASE: 0, ACK: 0, NACK: 0}


// not the swarm but all the drones
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
	seqNum[REQUEST] += 1
	formPermGroup()
	for _, otherDroneId := range permissionGroup {
		log.Println("Request seqNum " + strconv.Itoa(seqNum[REQUEST]))
		reqMsg := MaekawaMessage{drone.ID, otherDroneId, REQUEST, path}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, reqMsg, seqNum[REQUEST])
	}
}

func release() {
	ackNo = 0
	seqNum[RELEASE] += 1
	for _, otherDroneId := range permissionGroup {
		log.Println("Release seqNum " + strconv.Itoa(seqNum[RELEASE]))
		rlsMsg := MaekawaMessage{drone.ID, otherDroneId, RELEASE, myPathLock}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, rlsMsg, seqNum[RELEASE])
	}
}

func ack(dest string) {
	seqNum[ACK] += 1
	log.Println("Ack seqNum " + strconv.Itoa(seqNum[ACK]))
	ackMsg := MaekawaMessage{drone.ID, dest, ACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, ackMsg, seqNum[ACK])
}

func nack(dest string) {
	seqNum[NACK] += 1
	log.Println("Nack seqNum " + strconv.Itoa(seqNum[NACK]))
	nackMsg := MaekawaMessage{drone.ID, dest, NACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, nackMsg, seqNum[NACK])
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
	} else {
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
		moveDrone(msg.Path.To, 20)
		release();
	}
}

func handleNack(msg MaekawaMessage) {

}

func isIntersect(path1 PathLock, path2 PathLock) bool {
	return (dist3D_Segment_to_Segment(path1, path2) < 2)
}
