package paxos

import (
	"log"
	"time"
)

type Network interface {
	send(m Message)
	recv(timeout time.Duration) (Message, bool)
}

type PaxosNetwork struct {
	recvQueues map[int]chan Message
}

func newPaxosNetwork(agents ...int) *PaxosNetwork {
	pn := &PaxosNetwork{
		recvQueues: make(map[int]chan Message, 0),
	}

	for _, a := range agents {
		pn.recvQueues[a] = make(chan Message, 1024)
	}
	return pn
}

func (pn *PaxosNetwork) agentNetwork(id int) *agentNetwork {
	return &agentNetwork{id: id, PaxosNetwork: pn}
}

func (pn *PaxosNetwork) send(m Message) {
	log.Printf("nt: send %+v", m)
	pn.recvQueues[m.To] <- m
}

func (pn *PaxosNetwork) empty() bool {
	var n int
	for i, q := range pn.recvQueues {
		log.Printf("nt: %d left %d", i, len(q))
		n += len(q)
	}
	return n == 0
}

func (pn *PaxosNetwork) recvFrom(from int, timeout time.Duration) (Message, bool) {
	select {
	case m := <-pn.recvQueues[from]:
		log.Printf("nt: recv %+v", m)
		return m, true
	case <-time.After(timeout):
		return Message{}, false
	}
}

type agentNetwork struct {
	id int
	*PaxosNetwork
}

func (an *agentNetwork) send(m Message) {
	an.PaxosNetwork.send(m)
}

func (an *agentNetwork) recv(timeout time.Duration) (Message, bool) {
	return an.recvFrom(an.id, timeout)
}
