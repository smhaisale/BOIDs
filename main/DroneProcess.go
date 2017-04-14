package main

import (
    "fmt"
    "net/http"
    "log"
    "math"
    "time"
    "strconv"
)

var drone Drone

func main() {

    fmt.Println("Provide drone ID and port: ")

    var droneId string
    var port string
    fmt.Scanf("%s %s", &droneId, &port)

    http.HandleFunc("/heartbeat", heartbeat)

    http.HandleFunc("/getDroneInfo", getDroneInfo)

    http.HandleFunc("/updateSwarmInfo", updateSwarmInfo)

    http.HandleFunc("/moveToPosition", moveToPosition)

    // Start the environment server on localhost port 18841 and log any errors
    log.Println("http server started on :" + port)
    err := http.ListenAndServe(":" + port, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

    drone = Drone{droneId, Position{0, 1, 2}, DroneType{"0", "normal", Dimensions{1, 2, 3}, Dimensions{1, 2, 3}, Speed{1, 2, 3}}, Speed{1, 2, 3}}

    for {
        var x, y, z, vx, vy, vz float64
        fmt.Println("Enter new coordinates: ")
        fmt.Scanf("%f %f %f %f %f %f", x, y, z, vx, vy, vz)
        log.Println("Scanned ", x, y, z, vx, vy, vz)
        moveDrone(Position{x, y, z}, Speed{vx, vy, vz})
    }
}

func moveDrone(newPos Position, speed Speed) {
    log.Println("Moving to ", newPos)
    tX := (newPos.X - drone.Pos.X) / speed.VX
    tY := (newPos.Y - drone.Pos.Y) / speed.VY
    tZ := (newPos.Z - drone.Pos.Z) / speed.VZ

    t := math.Max(tX, math.Max(tY, tZ))

    for i := 0; i < int(t + 0.5); i++ {
        drone.Pos.X += (newPos.X - drone.Pos.X) / t
        drone.Pos.Y += (newPos.Y - drone.Pos.Y) / t
        drone.Pos.Z += (newPos.Z - drone.Pos.Z) / t
        time.Sleep(time.Duration(1000))
    }
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(drone.ID)))
}

func getDroneInfo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(drone)))
}

func updateSwarmInfo(w http.ResponseWriter, r *http.Request) {
    //w.Header().Set("Access-Control-Allow-Origin", "*")
    //w.Write([]byte(toJsonString(drone)))
}

func moveToPosition(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    x, _ := strconv.ParseFloat(values.Get("X"), 64)
    y, _ := strconv.ParseFloat(values.Get("Y"), 64)
    z, _ := strconv.ParseFloat(values.Get("Z"), 64)
    moveDrone(Position{x, y, z}, Speed{1, 1, 1})
}