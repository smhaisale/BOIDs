package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan UIMessage)         // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader {
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type UIMessage struct {
    MessageType     string  `json:"messageType"`
    DroneAddress    string  `json:"droneAddress"`
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