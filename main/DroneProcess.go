package main

import (
    "net/http"
    "strconv"
    "time"
    "log"
    "fmt"
    "math/rand"
    "math"
    "strings"
)

var droneObject DroneObject = DroneObject{}
var drone Drone = Drone{}
var swarm map[string]Drone = make(map[string]Drone)

var paxosClient = PaxosMessagePasser{}
var formPolygonPaxosClient = PaxosMessagePasser{}

var input_position = map[string]Position {
    "Drone0" : Position{0, 10, 0},
    "Drone1" : Position{0, 20, 0},
    "Drone2" : Position{0, 30, 0},
    "Drone3" : Position{0, 40, 0},
    "Drone4" : Position{0, 50, 0},
    "Drone5" : Position{0, 60, 0},
}

type MulticastMsgKey struct {
    OrigSender string
    Dest       string
    Type       string
    SeqNum     int
}

var haveSeenMap map[MulticastMsgKey]bool = make(map[MulticastMsgKey]bool)
var haveHandledMap map[MulticastMsgKey]bool = make(map[MulticastMsgKey]bool)
var pathLockManager = PathLockManager{permissionGroup: make([]string, 0), currPathLockList: make(map[string]PathLock), pathRequestQueue: make(map[string]PathLock), myPathLock: PathLock{}, ackNo: 0, seqNum: map[string]int{REQUEST: 0, RELEASE: 0, ACK: 0, NACK: 0}}


type MoveInstruction struct {
    Positions map[string]Position
}

func main() {

    rand.Seed( time.Now().UTC().UnixNano())

    var droneId, port string
    var x, y, z float64
    fmt.Println("Provide drone ID, port: ")
    fmt.Scanf("%s %s %f %f %f", &droneId, &port, x, y, z)

    paxosClient.id = 1
    formPolygonPaxosClient.id = 2

    http.HandleFunc(DRONE_HEARTBEAT_URL, heartbeat)
    http.HandleFunc(DRONE_GET_INFO_URL, getDroneInfo)
    http.HandleFunc(DRONE_UPDATE_SWARM_INFO_URL, updateSwarmInfo)
    http.HandleFunc(DRONE_MOVE_TO_POSITION_URL, moveToPosition)
    http.HandleFunc(DRONE_ADD_DRONE_URL, addNewDroneToSwarm)
    http.HandleFunc(DRONE_KILL_DRONE_URL, deleteDroneFromSwarm)
    http.HandleFunc(DRONE_PAXOS_MESSAGE_URL, handlePaxosMessage)
    http.HandleFunc(DRONE_FORM_POLYGON_URL, droneFormPolygon)
    http.HandleFunc(DRONE_FORM_SHAPE_URL, droneFormShape)
    http.HandleFunc("/proposeNewValue", proposeNewValue)
    http.HandleFunc(DRONE_MAEKAWA_MESSAGE_URL, handleMaekawaMessage)

    randomPosition := Position{rand.Float64() * 30 - 15, rand.Float64() * 20, rand.Float64() * 30 - 15}
    randomSpeed := Speed{rand.Float64() * 5, rand.Float64() * 5, rand.Float64() * 5}
    randomColor := Position{rand.Float64() * 255, rand.Float64() * 255, rand.Float64() * 255}

    droneObject = DroneObject{randomPosition, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, randomSpeed, randomColor, 0.7 + rand.Float64() * 0.3}
    drone = Drone{droneId, "localhost:" + port, droneObject}
    // Start the environment server and log any errors
    log.Println("http server started on " + drone.Address)
    err := http.ListenAndServe(":" + port, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
func move(newPos Position) {
    log.Println("Moving to ", newPos)
    oldPos := droneObject.Pos
    var deltaX, deltaY, deltaZ float64
    deltaX = newPos.X - oldPos.X
    deltaY = newPos.Y - oldPos.Y
    deltaZ = newPos.Z - oldPos.Z

    iterations := (math.Abs(deltaX) + math.Abs(deltaY) + math.Abs(deltaZ)) * 1.0

    for i := 0; i < int(iterations); i++ {
        droneObject.Pos.X += (deltaX) / iterations
        droneObject.Pos.Y += (deltaY) / iterations
        droneObject.Pos.Z += (deltaZ) / iterations

        time.Sleep(time.Duration(100000000))
        drone.DroneObject = droneObject
    }

    oldPosAfter := droneObject.Pos
    deltaX = newPos.X - oldPosAfter.X
    deltaY = newPos.Y - oldPosAfter.Y
    deltaZ = newPos.Z - oldPosAfter.Z

    droneObject.Pos.X += (deltaX)
    droneObject.Pos.Y += (deltaY)
    droneObject.Pos.Z += (deltaZ)

    time.Sleep(time.Duration(100000000))
    drone.DroneObject = droneObject

    log.Println("DroneObject in moveDrone", droneObject)
}

func moveDrone(newPos Position, t float64) {
    pathLockManager.request(PathLock{droneObject.Pos, newPos})
    //move(newPos)
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(drone.ID)))
}

func getDroneInfo(w http.ResponseWriter, r *http.Request) {
    drone.Address = r.Host
    // log.Println("Drone.droneObject in getDroneInfo ", drone.droneObject)
    // log.Println("DroneObject in moveDrone ", droneObject)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(drone)))
}

func updateSwarmInfo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(swarm)))
}

func getRandomPerpendicularPoints(from, to, between Position) (x, y, z float64) {
    x = between.X + 2 - rand.Float64() * 4
    y = between.Y + 2 - rand.Float64() * 4
    dX, dY, dZ := to.X - from.X, to.Y - from.Y, to.Z - from.Z
    z = (between.X * dX + between.Y * dY + between.Z * dZ - x * dX - y * dY) / dZ
    return
}

func moveToPosition(w http.ResponseWriter, r *http.Request) {

    values := r.URL.Query()
    x, _ := strconv.ParseFloat(values.Get("X"), 64)
    y, _ := strconv.ParseFloat(values.Get("Y"), 64)
    z, _ := strconv.ParseFloat(values.Get("Z"), 64)

	// todo make sure it works with Maekawa
	// Somehow ensure that the path is free - move other drones out of the way
    /*
    for _, swarmDrone := range swarm {
        flagY, flagZ := false, false
        deltaX := x - droneObject.Pos.X
        swarmDeltaX := x - swarmDrone.DroneObject.Pos.X
        if (z - droneObject.Pos.Z) / deltaX == (z - swarmDrone.DroneObject.Pos.Z) / swarmDeltaX {
            flagZ = true
        }
        if (y - droneObject.Pos.Y) / deltaX == (y - swarmDrone.DroneObject.Pos.Y) / swarmDeltaX {
            flagY = true
        }
        if flagY && flagZ {
            x, y, z := getRandomPerpendicularPoints(droneObject.Pos, Position{x, y, z}, swarmDrone.DroneObject.Pos)
            url := swarmDrone.Address + DRONE_MOVE_TO_POSITION_URL + "?X=" + strconv.FormatFloat(x, 'f', -1, 64) + "&Y=" + strconv.FormatFloat(y, 'f', -1, 64) + "&Z=" + strconv.FormatFloat(z, 'f', -1, 64)
            makeGetRequest("http://" + url, "")
        }
    }
    */

    moveDrone(Position{x, y, z}, 20)
}

func addNewDroneToSwarm(w http.ResponseWriter, r *http.Request) {
    address := r.URL.Query() .Get("address")
    log.Println("Received add drone request at address " + address)
    if address == drone.Address {
        return
    }
    newDrone, err := getDroneFromServer(address)
    if err != nil || swarm[newDrone.ID] == newDrone {
        log.Println("Error! ", err)
        return
    } else {
        swarm[newDrone.ID] = newDrone
        for _, swarmDrone := range swarm {
            swarmDroneAddress := "http://" + swarmDrone.Address + DRONE_ADD_DRONE_URL + "?address=" + address
            makeGetRequest(swarmDroneAddress, "")
            makeGetRequest( "http://" + address + DRONE_ADD_DRONE_URL + "?address=" + swarmDrone.Address, "")
        }
        makeGetRequest( "http://" + address + DRONE_ADD_DRONE_URL + "?address=" + drone.Address, "")
    }
}

func deleteDroneFromSwarm(w http.ResponseWriter, r *http.Request) {
    address := r.URL.Query() .Get("address")
    log.Println("Received kill drone request at address " + address)
    killDrone, err := getDroneFromServer(address)
    if err != nil {
        log.Println("Error! ", err)
        return
    } else {
        delete(swarm, killDrone.ID)
        for _, swarmDrone := range swarm {
            swarmDroneAddress := "http://" + swarmDrone.Address + DRONE_KILL_DRONE_URL + "?address=" + address
            makeGetRequest(swarmDroneAddress, "")
        }
    }
}

func proposeNewValue(w http.ResponseWriter, r *http.Request) {
    data := r.URL.Query() .Get("data")
    message := paxosClient.createPrepareMessage(data)
    paxosClient.sendPaxosMessage(message)
}

func handlePaxosMessage(w http.ResponseWriter, r *http.Request) {
    message := PaxosMessage{}
    getRequestBody(&message, r)

    switch message.ID {
    case 1:
        paxosClient.handlePaxosMessage(message)
    case 2:
        result := formPolygonPaxosClient.handlePaxosMessage(message)
        if result != "" {
            log.Println("Handle Paxos Message result : " + result)
            instruction := MoveInstruction{}
            fromJsonString(&instruction, result)
            moveDrone(instruction.Positions[drone.ID], 5)
        }
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func droneFormPolygon(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request at " + drone.ID)
    size, err := strconv.Atoi(r.URL.Query().Get("size"))
    if err != nil {
        size = len(swarm) + 1
    }
    index, positions := 0, getPolygonCoordinates(size, len(swarm) + 1)
    instruction := MoveInstruction{}
    instruction.Positions = map[string]Position{}
    for _, swarmDrone := range swarm {
        instruction.Positions[swarmDrone.ID] = positions[index]
        index++
    }
    instruction.Positions[drone.ID] = positions[index]
    message := formPolygonPaxosClient.createPrepareMessage(toJsonString(instruction))
    formPolygonPaxosClient.sendPaxosMessage(message)
}

func droneFormShape(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request at " + drone.ID)
    shape :=  r.URL.Query().Get("shape")
    size, err := strconv.Atoi(r.URL.Query().Get("size"))
    if err != nil {
        size = len(swarm) + 1
    }
    radius := 5 + rand.Float64() * 10
    positions := []Position{}
    if shape == "pyramid" {
        positions = calculateCoordinatesForPyramid(len(swarm) + 1, size,radius)
    } else if shape == "bipyramid" {
        positions = calculateCoordinatesForBipyramid(len(swarm) + 1, size,radius)
    } else if shape == "prism" {
        positions = calculateCoordinatesForPrism(len(swarm) + 1, size,radius)
    }
    index := 0
    instruction := MoveInstruction{}
    instruction.Positions = map[string]Position{}
    for _, swarmDrone := range swarm {
        instruction.Positions[swarmDrone.ID] = positions[index]
        index++
    }
    instruction.Positions[drone.ID] = positions[index]
    message := formPolygonPaxosClient.createPrepareMessage(toJsonString(instruction))
    formPolygonPaxosClient.sendPaxosMessage(message)
}

func handleMaekawaMessage(w http.ResponseWriter, r *http.Request) {
    origSender := r.URL.Query().Get("origSender")
    seqNum, _ := strconv.Atoi(r.URL.Query().Get("seqNum"))
    dest := r.URL.Query().Get("dest")
    msg := MaekawaMessage{}
    getRequestBody(&msg, r)
    //log.Println("Original " + msg.Type + " from " + origSender + " seq " + strconv.Itoa(seqNum))
    _, seen := haveSeenMap[MulticastMsgKey{origSender, dest, msg.Type, seqNum}]
    if !seen {
        haveSeenMap[MulticastMsgKey{origSender, dest, msg.Type, seqNum}] = true
        sendMulticast(DRONE_MAEKAWA_MESSAGE_URL, r.URL.Query(), msg)
    }
    _, handled := haveHandledMap[MulticastMsgKey{origSender, dest, msg.Type, seqNum}]
    if strings.Compare(drone.ID, dest) == 0 && !handled {
        haveHandledMap[MulticastMsgKey{origSender, dest, msg.Type, seqNum}] = true
        log.Println("Received Maekawa Message " + msg.Type + " from " + origSender + " with seq " + strconv.Itoa(seqNum))
        switch msg.Type {
        case REQUEST:
            pathLockManager.handleRequest(msg)
        case RELEASE:
            pathLockManager.handleRelease(msg)
        case ACK:
            pathLockManager.handleAck(msg)
        case NACK:
            handleNack(msg)
        }
    }
}
