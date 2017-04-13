package main

type Node struct {
	name    string
	address string
	port    string
}

var NodeMap = map[string]Node {
	"alice": Node{"alice", "localhost", "8081"},
	"bob": Node{"bob", "localhost", "8082"},
	"cat": Node{"cat", "localhost", "12358"},
	"deb": Node{"deb", "localhost", "12359"},
}
