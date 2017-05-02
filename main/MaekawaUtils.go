package main

import (
	"log"
	"strconv"
	"reflect"
	"math"
	"sync"
)

type PathLock struct {
	From Position
	To   Position
}

//var permissionGroup []string = make([]string, len(swarm))
//var currPathLockList map[string]PathLock = make(map[string]PathLock)
//var pathRequestQueue map[string]PathLock = make(map[string]PathLock)
//var myPathLock PathLock
//var ackNo int

type PathLockManager struct {
	permissionGroup []string
	currPathLockList map[string]PathLock
	pathRequestQueue map[string]PathLock
	myPathLock PathLock
	ackNo int
	seqNum map[string]int
	mux sync.Mutex
}

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

//var seqNum = map[string]int{REQUEST: 0, RELEASE: 0, ACK: 0, NACK: 0}


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
func (pathLockManager *PathLockManager) formPermGroup() {
	pathLockManager.permissionGroup = getDrones()
}

func (pathLockManager *PathLockManager) request(path PathLock) {
	pathLockManager.mux.Lock()
	pathLockManager.myPathLock = path
	pathLockManager.seqNum[REQUEST] += 1
	pathLockManager.formPermGroup()
	pathLockManager.mux.Unlock()
	for _, otherDroneId := range pathLockManager.permissionGroup {
		log.Println("Request seqNum " + strconv.Itoa(pathLockManager.seqNum[REQUEST]))
		reqMsg := MaekawaMessage{drone.ID, otherDroneId, REQUEST, path}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, reqMsg, pathLockManager.seqNum[REQUEST])
	}
}

func (pathLockManager *PathLockManager) release() {
	pathLockManager.mux.Lock()
	pathLockManager.ackNo = 0
	pathLockManager.seqNum[RELEASE] += 1
	pathLockManager.mux.Unlock()
	for _, otherDroneId := range pathLockManager.permissionGroup {
		log.Println("Release seqNum " + strconv.Itoa(pathLockManager.seqNum[RELEASE]))
		rlsMsg := MaekawaMessage{drone.ID, otherDroneId, RELEASE, pathLockManager.myPathLock}
		multicastMaekawa(drone.ID, otherDroneId, DRONE_MAEKAWA_MESSAGE_URL, rlsMsg, pathLockManager.seqNum[RELEASE])
	}
}

func (pathLockManager *PathLockManager) ack(dest string) {
	pathLockManager.mux.Lock()
	pathLockManager.seqNum[ACK] += 1
	pathLockManager.mux.Unlock()
	log.Println("Ack seqNum " + strconv.Itoa(pathLockManager.seqNum[ACK]))
	ackMsg := MaekawaMessage{drone.ID, dest, ACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, ackMsg, pathLockManager.seqNum[ACK])
}

func (pathLockManager *PathLockManager) nack(dest string) {
	pathLockManager.mux.Lock()
	pathLockManager.seqNum[NACK] += 1
	pathLockManager.mux.Unlock()
	log.Println("Nack seqNum " + strconv.Itoa(pathLockManager.seqNum[NACK]))
	nackMsg := MaekawaMessage{drone.ID, dest, NACK, PathLock{}}
	multicastMaekawa(drone.ID, dest, DRONE_MAEKAWA_MESSAGE_URL, nackMsg, pathLockManager.seqNum[NACK])
}

func (pathLockManager *PathLockManager) handleRequest(msg MaekawaMessage) {
	pathLockManager.mux.Lock()
	source := msg.Source
	path := msg.Path
	hasIntersect := false
	for _, pathLock := range pathLockManager.currPathLockList {
		if isIntersect(path, pathLock) {
			hasIntersect = true
			break
		}
	}
	if hasIntersect {
		pathLockManager.pathRequestQueue[source] = path
		//nack(source)
	} else {
		pathLockManager.currPathLockList[source] = path
		//ack(source)
	}
	pathLockManager.mux.Unlock()
	if hasIntersect {
		pathLockManager.nack(source)
	} else {
		pathLockManager.ack(source)
	}
}

func (pathLockManager *PathLockManager) handleRelease(msg MaekawaMessage) {
	pathLockManager.mux.Lock()
	source := msg.Source
	_, exist := pathLockManager.currPathLockList[source]
	if exist {
		delete(pathLockManager.currPathLockList, source)
	}
	var rlsSource []string
	for reqSource, reqPathLock := range pathLockManager.pathRequestQueue {
		hasIntersect := false
		for _, currPathLock := range pathLockManager.currPathLockList {
			if isIntersect(reqPathLock, currPathLock) {
				hasIntersect = true
				break
			}
		}
		if !hasIntersect {
			pathLockManager.currPathLockList[reqSource] = reqPathLock
			delete(pathLockManager.pathRequestQueue, reqSource)
			rlsSource = append(rlsSource, reqSource)
			//pathLockManager.ack(reqSource)
			break;
		}
	}
	pathLockManager.mux.Unlock()
	for _, src := range rlsSource {
		pathLockManager.ack(src)
	}
}

func (pathLockManager *PathLockManager) handleAck(msg MaekawaMessage) {
	pathLockManager.mux.Lock()
	pathLockManager.ackNo += 1
	if pathLockManager.ackNo >= len(pathLockManager.permissionGroup) {
		move(pathLockManager.myPathLock.To, 20)
		//moveDrone(pathLockManager.myPathLock.To, 20)
	}
	pathLockManager.mux.Unlock()
	if pathLockManager.ackNo >= len(pathLockManager.permissionGroup) {
		pathLockManager.release()
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
