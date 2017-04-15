package main

import (
    "log"
    "net/http"
    "math/rand"
)

// Get drone configuration from local cache instead of creating mock data.
var droneMap = map[string]Drone {}

type UIMessage struct {
    MessageType     string  `json:"messageType"`
    Data            string  `json:"data"`
}

func main() {

    http.HandleFunc(ENVIRONMENT_GET_ALL_DRONES_URL, getAllDrones)
    http.HandleFunc(ENVIRONMENT_ADD_DRONE_URL, addDrone)
    http.HandleFunc(ENVIRONMENT_KILL_DRONE_URL, killDrone)

    // Start the server on localhost port 8000 and log any errors
    log.Println("http server started on :18842")
    err := http.ListenAndServe(":18842", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func getRandomCoordinates () (x, y, z float64) {
    x = rand.Float64() * 20.0 - 10.0;
    y = rand.Float64() * 20.0;
    z = rand.Float64() * 20.0 - 10.0;
    log.Println("Random coordinates: ", x, y, z)
    return
}

func getAllDrones(w http.ResponseWriter, r *http.Request) {
    refreshDroneInfo()

    drones := []DroneObject{}
    for _, drone := range droneMap {
        drones = append(drones, drone.DroneObject)
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(drones)))
}

func addDrone(w http.ResponseWriter, r *http.Request) {

    address := r.URL.Query().Get("data")
    drone, err := getDroneFromServer(address)
    if err != nil {
        log.Println("Error! ", err)
    } else {
        droneMap[drone.ID] = Drone{drone.ID, address, "", drone.DroneObject}
    }
}

func killDrone(w http.ResponseWriter, r *http.Request) {

}

func refreshDroneInfo() {
    for key, drone := range droneMap {
        drone, err := getDroneFromServer(drone.Address)
        if err != nil {
            log.Println("Error in refreshDroneInfo()! ", err)
        } else {
            droneMap[key] = Drone{drone.ID, drone.Address,"", drone.DroneObject}
        }
    }
}