package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)        // connected clients
var broadcast = make(chan UIMessage)                // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader {
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type UIMessage struct {
    MessageType     string  `json:"messageType"`
    Data            string  `json:"data"`
}

func main() {

    //
    http.HandleFunc("/drones", handleDroneRequest)

    // Start listening for incoming chat messages
    go handleMessages()

    // Start the server on localhost port 8000 and log any errors
    log.Println("http server started on :18842")
    err := http.ListenAndServe(":18842", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleDroneRequest(w http.ResponseWriter, r *http.Request) {
    msg := new(UIMessage)
    getRequestBody(msg, r)
    log.Println(msg)

    // Get drone configuration from local cache instead of creating mock data.
    drones := []Drone {sampleDrone, sampleDrone, sampleDrone}
    drones[0].ID = "drone1"
    drones[0].Pos = Position {1, 1, 1}
    drones[1].ID = "drone2"
    drones[1].Pos = Position {2, 2, 2}
    drones[2].ID = "drone3"
    drones[2].Pos = Position {3, 3, 3}

    w.Write([]byte(toJsonString(drones)))
}


func handleMessages() {
    for {
        // Grab the next message from the broadcast channel
        msg := <-broadcast
        // Send it out to every client that is currently connected
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}