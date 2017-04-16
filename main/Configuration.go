package main

type Node struct {
	name    string
	address string
	port    string
}

var NodeMap = map[string]Node{
	"alice": Node{"alice", "localhost", "8081"},
	"bob":   Node{"bob", "localhost", "8082"},
	"cat":   Node{"cat", "localhost", "12358"},
	"deb":   Node{"deb", "localhost", "12359"},
}

// droneid: address
var DroneNodeMap = map[string]string {
	"drone1": "localhost:12345",
	"drone2": "localhost:12346",
	"drone3": "localhost:12347",
}

var MULTICAST_TYPE = "multicast"

var ENVIRONMENT_GET_ALL_DRONES_URL = "/getAllDrones"
var ENVIRONMENT_ADD_DRONE_URL = "/addDrone"
var ENVIRONMENT_KILL_DRONE_URL = "/killDrone"

var DRONE_ADD_DRONE_URL = "/killDrone"
var DRONE_GET_INFO_URL = "/getDroneInfo"
var DRONE_UPDATE_SWARM_INFO_URL = "/updateSwarmInfo"
var DRONE_MOVE_TO_POSITION_URL = "/moveToPosition"
var DRONE_HEARTBEAT_URL = "/heartbeat"

var DRONE_PAXOS_MESSAGE_URL = "/paxosMessage"
