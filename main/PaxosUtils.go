package main

import (
    "log"
)

var PAXOS_PREPARE_MESSAGE_TYPE = "prepare"
var PAXOS_PROMISE_MESSAGE_TYPE = "promise"
var PAXOS_ACCEPT_MESSAGE_TYPE = "accept"

var PAXOS_PROPOSER_ROLE_TYPE = "proposer"
var PAXOS_PROMISER_ROLE_TYPE = "promiser"
var PAXOS_ACCEPTOR_ROLE_TYPE = "acceptor"

type PaxosMessagePasser struct {
    id              int
    counter         map[string]int
    paxosSeqNum     int
    currentSource   string
    currentState    string
    currentValue    string
}

type PaxosMessage struct {
    ID              int
    Source          string
    Destination     string
    Type            string
    SeqNum          int
    Data            string
}

func (p *PaxosMessagePasser) createPrepareMessage(data string) (message PaxosMessage) {
    p.reset()
    p.currentState = PAXOS_PROPOSER_ROLE_TYPE
    message.Type = PAXOS_PREPARE_MESSAGE_TYPE
    p.paxosSeqNum++

    message.ID = p.id
    message.Source = drone.ID
    message.SeqNum = p.paxosSeqNum
    message.Data = data

    log.Println("Created prepare message: ", message)
    return message
}

func (p *PaxosMessagePasser) createPromiseMessage(dest string, data string) (message PaxosMessage) {
    message.ID = p.id
    message.Destination = dest
    message.Type = PAXOS_PROMISE_MESSAGE_TYPE
    message.SeqNum = p.paxosSeqNum
    message.Data = data

    log.Println("Created promise message: ", message)
    return message
}

func (p *PaxosMessagePasser) createAcceptMessage(data string) (message PaxosMessage) {

    message.ID = p.id
    message.Source = drone.ID
    message.Type = PAXOS_ACCEPT_MESSAGE_TYPE
    message.SeqNum = p.paxosSeqNum
    message.Data = data

    log.Println("Created accept message: ", message)
    return message
}

func (p *PaxosMessagePasser) sendPaxosMessage(message PaxosMessage) {
    jsonData := toJsonString(message)
    for _, drone := range swarm {
        address := "http://" + drone.Address + DRONE_PAXOS_MESSAGE_URL
        _, err := makeGetRequest(address, jsonData)
        if err != nil {
            log.Println("Error! ", err)
        }
    }
}

func (p *PaxosMessagePasser) handlePaxosMessage(message PaxosMessage) string {

    if message.SeqNum >= p.paxosSeqNum{
        log.Println("Received Paxos message: ", message)
    }

    // Based on current state and seq num, handle the message
    switch message.Type {
    case PAXOS_PREPARE_MESSAGE_TYPE:
        if message.SeqNum > p.paxosSeqNum ||
            (message.SeqNum == p.paxosSeqNum && message.Source > p.currentSource) {
            p.paxosSeqNum = message.SeqNum
            p.currentValue = message.Data
            p.currentState = PAXOS_PROMISER_ROLE_TYPE
            promiseMessage := p.createPromiseMessage(message.Source, p.currentValue)
            address := "http://" + swarm[message.Source].Address + DRONE_PAXOS_MESSAGE_URL
            _, err := makeGetRequest(address, toJsonString(promiseMessage))
            if err != nil {
                log.Println("Error! ", err)
            }
        }
    case PAXOS_PROMISE_MESSAGE_TYPE:
        log.Println(p)
        if p.currentState == PAXOS_PROPOSER_ROLE_TYPE {
            p.counter[message.Data]++
            log.Println(p.counter)
            log.Println(swarm)
            if p.counter[message.Data] >= len(swarm) / 2 && p.paxosSeqNum == message.SeqNum {
                p.currentValue = message.Data
                acceptMessage := p.createAcceptMessage(message.Data)
                p.sendPaxosMessage(acceptMessage)
                log.Printf("Setting accepted global value: ", message.Data)
                p.reset()
                return message.Data
            }
        }
    case PAXOS_ACCEPT_MESSAGE_TYPE:
        if p.paxosSeqNum <= message.SeqNum {
            p.paxosSeqNum = message.SeqNum
            log.Printf("Setting accepted global value: ", message.Data)
            p.reset()
            return message.Data
        }
    }
    return ""
}

func (p *PaxosMessagePasser) reset() {
    p.currentState = ""
    p.currentValue = ""
    p.currentSource = ""
    p.counter = make(map[string]int)
}