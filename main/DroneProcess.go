package main

import (
    "fmt"
    "net/http"
    "strconv"
    "time"
    "log"
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

type MoveInstruction struct {
    positions map[string]Position
}

func main() {

    var droneId, port string
    var x, y, z float64
    fmt.Println("Provide drone ID, port: ")
    fmt.Scanf("%s %s %f %f %f", &droneId, &port, x, y, z)

    http.HandleFunc(DRONE_HEARTBEAT_URL, heartbeat)
    http.HandleFunc(DRONE_GET_INFO_URL, getDroneInfo)
    http.HandleFunc(DRONE_UPDATE_SWARM_INFO_URL, updateSwarmInfo)
    http.HandleFunc(DRONE_MOVE_TO_POSITION_URL, moveToPosition)
    http.HandleFunc(DRONE_ADD_DRONE_URL, addNewDroneToSwarm)
    http.HandleFunc(DRONE_PAXOS_MESSAGE_URL, handlePaxosMessage)
    http.HandleFunc(DRONE_FORM_POLYGON_URL, droneFormPolygon)
    http.HandleFunc("/proposeNewValue", proposeNewValue)

    droneObject = DroneObject{Position{x, y, z}, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, Speed{1, 2, 3}}
    drone = Drone{droneId, "localhost:" + port, droneObject}
    // Start the environment server on localhost port 18841 and log any errors
    log.Println("http server started on :" + port)
    err := http.ListenAndServe(":" + port, nil)
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
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(swarm)))
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

func proposeNewValue(w http.ResponseWriter, r *http.Request) {
    data := r.URL.Query() .Get("data")
    message := paxosClient.createPrepareMessage(data)
    paxosClient.sendPaxosMessage(message)
}

func handlePaxosMessage(w http.ResponseWriter, r *http.Request) {
    message := PaxosMessage{}
    getRequestBody(&message, r)

    switch (message.ID) {
    case 1:
        paxosClient.handlePaxosMessage(message)
    case 2:
        result := formPolygonPaxosClient.handlePaxosMessage(message)
        if result != "" {
            instruction := MoveInstruction{}
            fromJsonString(instruction, result)
            moveDrone(instruction.positions[drone.ID], 10)
        }
    }
}

func droneFormPolygon(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request at " + drone.ID)
    index, positions := 0, calculateCoordinates(len(swarm))
    instruction := MoveInstruction{}
    for _,drone := range swarm {
        instruction.positions[drone.ID] = positions[index]
        index++
    }
    message := formPolygonPaxosClient.createPrepareMessage(toJsonString(instruction))
    formPolygonPaxosClient.sendPaxosMessage(message)
}