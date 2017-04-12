package main

type Node struct {
	name    string
	address string
	port    string
}

var NodeMap = map[string]Node{
	"alice": Node{"alice", "192.168.1.5", "8081"},
	"bob": Node{"bob", "192.168.1.5", "8082"},
	"cat": Node{"cat", "192.168.1.5", "12358"},
	"deb": Node{"deb", "192.168.1.5", "12359"},
}
