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
    http.HandleFunc(ENVIRONMENT_FORM_POLYGON_URL, formPolygon)

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

    drones := []Drone{}
    for _, drone := range droneMap {
        drones = append(drones, drone)
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(toJsonString(drones)))
}

func addDrone(w http.ResponseWriter, r *http.Request) {

    address := r.URL.Query().Get("data")
    log.Println("Received add drone request at address " + address)
    drone, err := getDroneFromServer(address)
    if err != nil {
        log.Println("Error! ", err)
    } else {
        droneMap[drone.ID] = Drone{drone.ID, address, drone.DroneObject}
        for _, swarmDrone := range droneMap {
            swarmDroneAddress := "http://" + swarmDrone.Address + DRONE_ADD_DRONE_URL + "?address=" + address
            makeGetRequest(swarmDroneAddress, "")
        }
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func killDrone(w http.ResponseWriter, r *http.Request) {

}

func formPolygon(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request")
    for _, drone := range droneMap {
        address := "http://" + drone.Address + ENVIRONMENT_FORM_POLYGON_URL
        makeGetRequest(address, "")
        break
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func refreshDroneInfo() {
    for key, drone := range droneMap {
        drone, err := getDroneFromServer(drone.Address)
        if err != nil {
            log.Println("Error in refreshDroneInfo()! ", err)
        } else {
            droneMap[key] = Drone{drone.ID, drone.Address, drone.DroneObject}
        }
    }
}