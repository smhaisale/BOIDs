package main

// Message structure used for TCP connections between drones
// Contains message metadata and payload
type TcpMessage struct {
    Source          string          `json: "source"`
    Destination     string          `json: "destination"`
    SeqNum          int             `json: "seqNum"`
    Duplicate       bool            `json: "duplicate"`
    Type            string          `json: "type"`
    Data            MessageData     `json: "data"`
    MulticastData   MulticastData   `json: "multicastData"`
}

// Contains message payload
type MessageData struct {
    Drones          []Drone         `json: "drones"`
}

// Contains metadata specific to multicast messages
type MulticastData struct {
    Source          string          `json: "source"`
    Destination     string          `json: "destination"`
    GroupSeqNum     int             `json: "groupSeqNum"`
    GroupName       string          `json: "groupName"`
}