package paxos

import (
	"log"
	"time"
)

type Learner struct {
	id        int
	acceptors map[int]accept

	nt Network
}

func newLearner(id int, nt Network, acceptors ...int) *Learner {
	l := &Learner{id: id, nt: nt, acceptors: make(map[int]accept)}
	for _, a := range acceptors {
		l.acceptors[a] = Message{Type: Accept}
	}
	return l
}

// A Value is learned when a single proposal with that Value has been accepted by
// a majority of the acceptors.
func (l *Learner) learn() string {
	for {
		m, ok := l.nt.recv(time.Hour)
		if !ok {
			continue
		}
		if m.Type != Accept {
			log.Panicf("Learner: %d received unexpected msg %+v", l.id, m)
		}
		l.receiveAccepted(m)
		accept, ok := l.chosen()
		if !ok {
			continue
		}
		log.Printf("Learner: %d learned the chosen propose %+v", l.id, accept)
		return accept.proposalValue()
	}
}

func (l *Learner) receiveAccepted(accepted Message) {
	a := l.acceptors[accepted.From]
	if a.proposalNumber() < accepted.N {
		log.Printf("Learner: %d received a new accepted proposal %+v", l.id, accepted)
		l.acceptors[accepted.From] = accepted
	}
}

func (l *Learner) majority() int { return len(l.acceptors)/2 + 1 }

// A proposal is chosen when it has been accepted by a majority of the
// acceptors.
// The leader might choose multiple proposals when it learns multiple times,
// but we guarantee that all chosen proposals have the same Value.
func (l *Learner) chosen() (accept, bool) {
	counts := make(map[int]int)
	accepteds := make(map[int]accept)

	for _, accepted := range l.acceptors {
		if accepted.proposalNumber() != 0 {
			counts[accepted.proposalNumber()]++
			accepteds[accepted.proposalNumber()] = accepted
		}
	}

	for n, count := range counts {
		if count >= l.majority() {
			return accepteds[n], true
		}
	}
	return Message{}, false
}
