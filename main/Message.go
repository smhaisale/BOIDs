package main

// Message structure used for TCP connections between drones
// Contains message metadata and payload
type TcpMessage struct {
	source        string        `json: "source"`
	destination   string        `json: "destination"`
	seqNum        int           `json: "seqNum"`
	duplicate     bool          `json: "duplicate"`
	msgType       string        `json: "type"`
	timestamp     VectorTime    `json: "vectorTimestamp"`
	data          MessageData   `json: "data"`
	multicastData MulticastData `json: "multicastData"`
}

// Contains message payload
type MessageData struct {
	drones []Drone `json: "drones"`
}

// Contains metadata specific to multicast messages
type MulticastData struct {
	source      string `json: "source"`
	destination string `json: "destination"`
	groupSeqNum int    `json: "groupSeqNum"`
	groupName   string `json: "groupName"`
}
