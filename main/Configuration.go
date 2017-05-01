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

var ENVIRONMENT_GET_ALL_DRONES_URL = "/getAllDrones"
var ENVIRONMENT_ADD_DRONE_URL = "/addDrone"
var ENVIRONMENT_KILL_DRONE_URL = "/killDrone"
var ENVIRONMENT_FORM_POLYGON_URL = "/formPolygon"
var ENVIRONMENT_FORM_SHAPE_URL = "/formShape"
var ENVIRONMENT_RANDOM_POSITIONS_URL = "/randomPositions"

var DRONE_ADD_DRONE_URL = "/addDroneToSwarm"
var DRONE_KILL_DRONE_URL = "/deleteDroneFromSwarm"
var DRONE_GET_INFO_URL = "/getDroneInfo"
var DRONE_UPDATE_SWARM_INFO_URL = "/getSwarmInfo"
var DRONE_MOVE_TO_POSITION_URL = "/moveToPosition"
var DRONE_HEARTBEAT_URL = "/heartbeat"
var DRONE_FORM_POLYGON_URL = "/formPolygon"
var DRONE_FORM_SHAPE_URL = "/formShape"

var DRONE_PAXOS_MESSAGE_URL = "/paxosMessage"

var DRONE_MAEKAWA_MESSAGE_URL = "/maekawaMessage"