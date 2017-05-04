package main

import (
    "log"
    "net/http"
    "math/rand"
    "strconv"
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
    http.HandleFunc(ENVIRONMENT_FORM_SHAPE_URL, formShape)
    http.HandleFunc(ENVIRONMENT_RANDOM_POSITIONS_URL, randomPositions)

    // Start the server on localhost port 8000 and log any errors
    log.Println("http server started on :18842")
    err := http.ListenAndServe(":18842", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func getRandomCoordinates () (x, y, z string) {
    x2 := rand.Float64() * 30.0 - 15.0
    y2 := rand.Float64() * 20.0
    z2 := rand.Float64() * 30.0 - 15.0
    log.Println("Random coordinates: ", x, y, z)
    x = strconv.FormatFloat(x2, 'f', 6, 64)
    y = strconv.FormatFloat(y2, 'f', 6, 64)
    z = strconv.FormatFloat(z2, 'f', 6, 64)
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

func kill(droneId string) {
    log.Println("Received kill drone request for " + droneId)
    address := droneMap[droneId].Address
    killDrone, _ := getDroneFromServer(address)
    delete(droneMap, killDrone.ID)
    for _, swarmDrone := range droneMap {
        killDroneAddress := "http://" + swarmDrone.Address + DRONE_KILL_DRONE_URL + "?address=" + address
        makeGetRequest(killDroneAddress, "")
    }
}

func killDrone(w http.ResponseWriter, r *http.Request) {
    droneId := r.URL.Query().Get("data")
    kill(droneId)
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func formPolygon(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request")
    size := r.URL.Query().Get("nodes")
    for _, drone := range droneMap {
        address := "http://" + drone.Address + DRONE_FORM_POLYGON_URL + "?size=" + size
        asyncGetRequest(address, "")
        break
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func formShape(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request")
    shape := r.URL.Query().Get("shape")
    size := r.URL.Query().Get("nodes")
    for _, drone := range droneMap {
        address := "http://" + drone.Address + DRONE_FORM_SHAPE_URL + "?shape=" + shape + "&size=" + size
        asyncGetRequest(address, "")
        break
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func randomPositions(w http.ResponseWriter, r *http.Request) {
    log.Println("Received form polygon request")
    for _, drone := range droneMap {
        address := "http://" + drone.Address + DRONE_MOVE_TO_POSITION_URL
        x, y, z := getRandomCoordinates()
        asyncGetRequest(address + "?X=" + x + "&Y=" + y + "&Z=" + z, "")
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func refreshDroneInfo() {
    for key, drone := range droneMap {
        drone, err := getDroneFromServer(drone.Address)
        if err != nil {
            kill(key)
            log.Println("Error in refreshDroneInfo()! ", err)
        } else {
            droneMap[key] = Drone{drone.ID, drone.Address, drone.DroneObject}
        }
    }
}
