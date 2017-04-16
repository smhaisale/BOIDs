package main

import (
    "net/http"
    "log"
)

type PathLock struct {
    From    Position
    To      Position
}

var permissionGroup []string

var locks []PathLock

type MaekawaMessage struct {
    Source          string
    Destination     string
    Type            string
    SeqNum          int
    Data            string
}

func createLockMessage(data string) (message MaekawaMessage) {
    reset()
    message.Source = drone.ID
    message.Data = data

    log.Println("Created prepare message: ", message)
    return message
}

func createReleaseMessage(dest string, data string) (message MaekawaMessage) {
    message.Destination = dest
    message.Data = data

    log.Println("Created promise message: ", message)
    return message
}

func createRejectMessage(data string) (message MaekawaMessage) {
    message.Source = drone.ID
    message.Data = data

    log.Println("Created accept message: ", message)
    return message
}

func createFailedMessage(message MaekawaMessage) {

}

func handleMaekawaMessage(w http.ResponseWriter, r *http.Request) {
    message := MaekawaMessage{}
    getRequestBody(&message, r)
}