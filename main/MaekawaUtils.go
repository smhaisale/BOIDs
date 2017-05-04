package main

import (
	"log"
	"strconv"
	"reflect"
	"math"
)

type PathLock struct {
	From Position
	To   Position
}

var permissionGroup []string = make([]string, 0)
var currPathLockList map[string]PathLock = make(map[string]PathLock)
var pathRequestQueue map[string]PathLock = make(map[string]PathLock)
var myPathLock = PathLock{}
var ackNo int = 0
var seqNum = map[string]int{REQUEST:0, RELEASE:0, ACK:0, NACK:0}
/*
type PathLockManager struct {
	permissionGroup []string
	currPathLockList map[string]PathLock
	pathRequestQueue map[string]PathLock
	myPathLock PathLock
	ackNo int
	seqNum map[string]int
	mux sync.Mutex
}
*/

type MaekawaMessage struct {
	Source      string	`json: "source"`
	Destination string	`json: "dest"`
	Type        string	`json: "type"`
	Path        PathLock	`json: "path"`
}

var REQUEST = "REQUEST"
var RELEASE = "RELEASE"
var ACK = "ACK"
var NACK = "NACK"


// not the swarm but all the drones
func getDrones() []string{
	keys := reflect.ValueOf(swarm).MapKeys()
	drones := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		drones[i] = keys[i].String()
	}
	drones = append(drones, drone.ID)
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
	formPermGroup()
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
		//move(pathLockManager.myPathLock.To, 5)
		move(myPathLock.To)
		release()
	}

}

func handleNack(msg MaekawaMessage) {

}

func isIntersect(path1 PathLock, path2 PathLock) bool {
	return (dist3D_Segment_to_Segment(path1, path2) <= 0)
}

func dist3D_Segment_to_Segment(path1 PathLock, path2 PathLock) float64 {
	SMALL_NUM := 0.00000001
	u := sub(path1.To, path1.From)
	v := sub(path2.To, path2.From)
	w := sub(path1.From, path2.From)
	a := dotMul(u, u)
	b := dotMul(u, v)
	c := dotMul(v, v)
	d := dotMul(u, w)
	e := dotMul(v, w)
	D := a * c - b * b
	sc, sN, sD := D, D, D
	tc, tN, tD := D, D, D

	if D < SMALL_NUM {
		sN = 0.0
		sD = 1.0
		tN = e
		tD = c
	} else {
		sN = b * e - c * d
		tN = a * e - b * d
		if sN < 0.0 {
			sN = 0.0
			tN = e
			tD = c
		} else if sN > sD {
			sN = sD
			tN = e + b
			tD = c
		}
	}
	if tN < 0.0 {
		tN = 0.0
		if -d < 0.0 {
			sN = 0.0
		} else if -d > a {
			sN = sD
		} else {
			sN = -d
			sD = a
		}
	} else if tN > tD {
		tN = tD
		if -d + b < 0.0 {
			sN = 0
		} else if -d + b > a {
			sN = sD
		} else {
			sN = -d +  b
			sD = a
		}
	}
	if math.Abs(sN) < SMALL_NUM {
		sc = 0.0
	} else {
		sc = sN / sD
	}
	if math.Abs(tN) < SMALL_NUM {
		tc = 0.0
	} else {
		tc = tN / tD
	}

	// get the difference of the two closest points
	dP := sub(add(w, scalaMul(sc, u)), scalaMul(tc, v))

	return norm(dP);   // return the closest distance

}
