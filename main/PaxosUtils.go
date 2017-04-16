package main

type PaxosMessage struct {
	from, to    string
	messageType string
	value       string
}
