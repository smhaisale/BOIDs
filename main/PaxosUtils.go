package main

import (
    "net/http"
    "log"
)

var PAXOS_PREPARE_MESSAGE_TYPE = "prepare"
var PAXOS_PROMISE_MESSAGE_TYPE = "promise"
var PAXOS_ACCEPT_MESSAGE_TYPE = "accept"

var PAXOS_PROPOSER_ROLE_TYPE = "proposer"
var PAXOS_PROMISER_ROLE_TYPE = "promiser"
var PAXOS_ACCEPTOR_ROLE_TYPE = "acceptor"


var counter map[string]int

var paxosSeqNum = 0
var currentPromisedNumber = 0
var currentAcceptedNumber = 0

var currentSource = ""
var currentState = ""
var currentValue = ""
var currentTestValue = ""

type PaxosMessage struct {
    Source          string
    Destination     string
    Type            string
    SeqNum          int
    Data            string
}

func createPrepareMessage(data string) (message PaxosMessage) {
    reset()
    currentState = PAXOS_PROPOSER_ROLE_TYPE
    message.Source = drone.ID
    message.Type = PAXOS_PREPARE_MESSAGE_TYPE
    paxosSeqNum++
    message.SeqNum = paxosSeqNum
    message.Data = data

    log.Println("Created prepare message: ", message)
    return message
}

func createPromiseMessage(dest string, data string) (message PaxosMessage) {
    message.Destination = dest
    message.Type = PAXOS_PROMISE_MESSAGE_TYPE
    message.SeqNum = paxosSeqNum
    message.Data = data

    log.Println("Created promise message: ", message)
    return message
}

func createAcceptMessage(data string) (message PaxosMessage) {
    message.Source = drone.ID
    message.Type = PAXOS_ACCEPT_MESSAGE_TYPE
    message.SeqNum = paxosSeqNum
    message.Data = data

    log.Println("Created accept message: ", message)
    return message
}

func sendPaxosMessage(message PaxosMessage) {
    jsonData := toJsonString(message)
    for _, drone := range swarm {
        address := "http://" + drone.Address + DRONE_PAXOS_MESSAGE_URL
        _, err := makeGetRequest(address, jsonData)
        if err != nil {
            log.Println("Error! ", err)
        }
    }
}

func handlePaxosMessage(w http.ResponseWriter, r *http.Request) {
    message := PaxosMessage{}
    getRequestBody(&message, r)

    if message.SeqNum >= paxosSeqNum{
        log.Println("Received Paxos message: ", message)
    }

    // Based on current state and seq num, handle the message
    switch message.Type {
    case PAXOS_PREPARE_MESSAGE_TYPE:
        if message.SeqNum > paxosSeqNum ||
            (message.SeqNum == paxosSeqNum && message.Source > currentSource) {
            paxosSeqNum = message.SeqNum
            currentValue = message.Data
            currentState = PAXOS_PROMISER_ROLE_TYPE
            promiseMessage := createPromiseMessage(message.Source, currentValue)
            address := "http://" + swarm[message.Source].Address + DRONE_PAXOS_MESSAGE_URL
            _, err := makeGetRequest(address, toJsonString(promiseMessage))
            if err != nil {
                log.Println("Error! ", err)
            }
        }
    case PAXOS_PROMISE_MESSAGE_TYPE:
        if currentState == PAXOS_PROPOSER_ROLE_TYPE && paxosSeqNum == message.SeqNum {
            counter[message.Data]++
            if counter[message.Data] > len(swarm) / 2 {
                currentValue = message.Data
                // Set accepted value to above
                acceptMessage := createAcceptMessage(currentValue)
                sendPaxosMessage(acceptMessage)
                log.Printf("Setting accepted global value: ", message.Data)
                reset()
            }
        }
    case PAXOS_ACCEPT_MESSAGE_TYPE:
        if paxosSeqNum <= message.SeqNum {
            paxosSeqNum = message.SeqNum
            log.Printf("Setting accepted global value: ", message.Data)
            // Set accepted values according to message
            reset()
        }
    }
}

func reset() {
    currentState = ""
    currentValue = ""
    currentSource = ""
    counter = make(map[string]int)
}