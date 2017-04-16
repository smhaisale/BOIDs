package main

import (
	"./paxos"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"net/url"
)

var droneObject DroneObject
var drone Drone
var swarm = make(map[string]Drone)

type MulticastRequestKey struct {
	OriginalSender string
	GroupSeqNum    int
}
var haveSeenMap = make(map[MulticastRequestKey]bool)

var input_position = map[string]Position{
	"Drone0": Position{0, 10, 0},
	"Drone1": Position{0, 20, 0},
	"Drone2": Position{0, 30, 0},
	"Drone3": Position{0, 40, 0},
	"Drone4": Position{0, 50, 0},
	"Drone5": Position{0, 60, 0},
}

func main() {

	var droneId, port, paxosRole string
	fmt.Println("Provide drone ID, port, paxosRole: ")
	fmt.Scanf("%s %s", &droneId, &port, &paxosRole)

	http.HandleFunc(DRONE_HEARTBEAT_URL, heartbeat)
	http.HandleFunc(DRONE_GET_INFO_URL, getDroneInfo)
	http.HandleFunc(DRONE_UPDATE_SWARM_INFO_URL, updateSwarmInfo)
	http.HandleFunc(DRONE_MOVE_TO_POSITION_URL, moveToPosition)
	http.HandleFunc(DRONE_ADD_DRONE_URL, addNewDroneToSwarm)
	http.HandleFunc(DRONE_PAXOS_MESSAGE_URL, handlePaxosMessage)
	http.HandleFunc("/getswarm", getSwarm)
	http.HandleFunc("/multi", multi)

	droneObject = DroneObject{Position{0, 0, 0}, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, Speed{1, 2, 3}}
	drone = Drone{droneId, "localhost:" + port, paxosRole, droneObject}
	// Start the environment server on localhost port 18841 and log any errors
	swarm[droneId] = drone
	log.Println("http server started on :" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//func moveDrone(newPos Position, speed Speed) {
//    log.Println("Moving to ", newPos)
//    tX := math.Abs((newPos.X - drone.Pos.X) / speed.VX)
//    tY := math.Abs((newPos.Y - drone.Pos.Y) / speed.VY)
//    tZ := math.Abs((newPos.Z - drone.Pos.Z) / speed.VZ)
//
//    t := math.Max(tX, math.Max(tY, tZ))
//
//    for i := 0; i < int(t + 0.5); i++ {
//        drone.Pos.X += (newPos.X - drone.Pos.X) / t
//        drone.Pos.Y += (newPos.Y - drone.Pos.Y) / t
//        drone.Pos.Z += (newPos.Z - drone.Pos.Z) / t
//        time.Sleep(time.Duration(1000000000))
//    }
//}

func moveDrone(newPos Position, t float64) {
	log.Println("Moving to ", newPos)
	oldPos := droneObject.Pos
	for {
		if newPos.X == droneObject.Pos.X && newPos.Y == droneObject.Pos.Y && newPos.Z == droneObject.Pos.Z {
			break
		}
		droneObject.Pos.X += (newPos.X - oldPos.X) / t
		droneObject.Pos.Y += (newPos.Y - oldPos.Y) / t
		droneObject.Pos.Z += (newPos.Z - oldPos.Z) / t
		time.Sleep(time.Duration(1000000000))
		drone.DroneObject = droneObject
	}
	log.Println("DroneObject in moveDrone", droneObject)

}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(toJsonString(drone.ID)))
}

func getDroneInfo(w http.ResponseWriter, r *http.Request) {
	// log.Println("Drone.droneObject in getDroneInfo ", drone.droneObject)
	// log.Println("DroneObject in moveDrone ", droneObject)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(toJsonString(drone)))
}

func updateSwarmInfo(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Write([]byte(toJsonString(drone)))
}

func getSwarm(w http.ResponseWriter, r *http.Request) {
	drones := getDrones()
	w.Write([]byte(toJsonString(drones)))
}

func multi(w http.ResponseWriter, r *http.Request) {
	port := r.URL.Query().Get("port")
	param := url.Values{}
	param.Add("type", "multicast")
	param.Add("address", "localhost:"+port)
	Multicast("drone1", "drone2", DRONE_ADD_DRONE_URL, param, "")
}

func moveToPosition(w http.ResponseWriter, r *http.Request) {
	// log.Println("Drone.droneObject in moveToPosition ", drone.droneObject)
	//  log.Println("DroneObject in moveToPosition ", droneObject)
	values := r.URL.Query()
	x, _ := strconv.ParseFloat(values.Get("X"), 64)
	y, _ := strconv.ParseFloat(values.Get("Y"), 64)
	z, _ := strconv.ParseFloat(values.Get("Z"), 64)
	moveDrone(Position{x, y, z}, 20)
}

func addNewDroneToSwarm(w http.ResponseWriter, r *http.Request) {
	msgType := r.URL.Query().Get("type")
	if strings.Compare(msgType, MULTICAST_TYPE) == 0 {
		msg := MulticastMessage{}
		getRequestBody(&msg, r)
		seenMsgKey := MulticastRequestKey{msg.OriginalSender, msg.GroupSeqNum}
		if !haveSeenMap[seenMsgKey] {
			address := r.URL.Query().Get("address")
			newDrone, err := getDroneFromServer(address)
			if err != nil {
				log.Println("Error! ", err)
			} else {
				swarm[newDrone.ID] = newDrone
			}
			haveSeenMap[seenMsgKey] = true
			defer sendMulticast(DRONE_ADD_DRONE_URL, r.URL.Query(), msg)
		}
	} else {
		address := r.URL.Query().Get("address")
		newDrone, err := getDroneFromServer(address)
		if err != nil {
			log.Println("Error! ", err)
		} else {
			swarm[newDrone.ID] = newDrone
		}
	}
	drones := getDrones()
	w.Write([]byte(toJsonString(drones)))
}

func sendPaxosMessage(droneId string, paxosMessage paxos.Message) {

}

func handlePaxosMessage(w http.ResponseWriter, r *http.Request) {

}
