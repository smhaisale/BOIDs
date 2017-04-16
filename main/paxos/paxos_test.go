package paxos

import (
	"testing"
	"time"
)

func TestPaxosWithSingleProposer(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001 is a Proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 2001)

	as := make([]*Acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p := newProposer(1001, "3D coordinates", pn.agentNetwork(1001), 1, 2, 3)
	go p.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "3D coordinates" {
		t.Errorf("Value = %s, want %s", value, "3D coordinates")
	}
}

func TestPaxosWithTwoProposers(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001,1002 is a Proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 2001)

	as := make([]*Acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p1 := newProposer(1001, "3D coordinates", pn.agentNetwork(1001), 1, 2, 3)
	go p1.run()

	time.Sleep(time.Millisecond)
	p2 := newProposer(1002, "4D coordinates", pn.agentNetwork(1002), 1, 2, 3)
	go p2.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "3D coordinates" {
		t.Errorf("Value = %s, want %s", value, "3D coordinates")
	}
	time.Sleep(time.Millisecond)
}
