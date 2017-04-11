package main

import (
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
    "fmt"
)

// Message structure used for TCP connections between drones
// Contains list of drones
type TcpMessage struct {
    Drones      []Drone       `json: "drone"`
}

// All contained variable names must begin with a capital letter to be visible by JSONWrapper
func toJsonString(message TcpMessage) string {
    msg, err := json.Marshal(message)
    if err != nil {
        log.Fatal(err)
    }
    return string(msg)
}

func main() {
    var t = TcpMessage{[]Drone{sampleDrone, sampleDrone, sampleDrone}}
    fmt.Println(toJsonString(t))
}

func getRequestBody(msg interface {}, req *http.Request) interface{} {

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))
    err = json.Unmarshal(body, msg)
    if err != nil {
        log.Printf("error: %v", err)
    }
    defer req.Body.Close()
    return msg
}
