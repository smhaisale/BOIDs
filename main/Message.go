package main

// Message structure used for TCP connections between drones
// Contains message metadata and payload
type TcpMessage struct {
	Source        string        `json: "source"`
	Destination   string        `json: "destination"`
	SeqNum        int           `json: "seqNum"`
	Duplicate     bool          `json: "duplicate"`
	MsgType       string        `json: "msgType"`
	Timestamp     VectorTime    `json: "vectorTimestamp"`
	Data          MessageData   `json: "data"`
	MulticastData MulticastData `json: "multicastData"`
}

// Contains message payload
type MessageData struct {
	Drones []Drone `json: "drones"`
}

// Contains metadata specific to multicast messages
type MulticastData struct {
	Source      string `json: "source"`
	Destination string `json: "destination"`
	GroupSeqNum int    `json: "groupSeqNum"`
	GroupName   string `json: "groupName"`
}

type MulticastMessage struct {
	OriginalSender string `json: "originalSender"`
	Destination    string `json: "destination"`
	GroupSeqNum    int    `json: "groupSeqNum"`
	MessageData    string `json: messageData`
}

var sampleMulticastData = MulticastData{"multicastSource", "multicastDestination",
	1, "groupName"}

var sampleMessageData = MessageData{[]Drone{sampleDrone, sampleDrone, sampleDrone}}

var sampleTcpMessage = TcpMessage{"source", "destination", 1, false, "messageType",
	sampleTimestamp, sampleMessageData, sampleMulticastData}
