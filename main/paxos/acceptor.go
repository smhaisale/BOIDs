package paxos

import (
	"log"
	"time"
)

// P1. An Acceptor must accept the first proposal that it receives.
// If a proposal with Value v is chosen, then every higher-numbered proposal
// accepted by any Acceptor has Value v.
type Acceptor struct {
	id       int
	learners []int

	accept   Message
	promised promise

	nt Network
}

func newAcceptor(id int, nt Network, learners ...int) *Acceptor {
	return &Acceptor{id: id, nt: nt, promised: Message{}, learners: learners}
}

func (a *Acceptor) run() {
	for {
		m, ok := a.nt.recv(time.Hour)
		if !ok {
			continue
		}
		switch m.Type {
		case Propose:
			accepted := a.receivePropose(m)
			if accepted {
				for _, l := range a.learners {
					m := a.accept
					m.From = a.id
					m.To = l
					a.nt.send(m)
				}
			}
		case Prepare:
			promise, ok := a.receivePrepare(m)
			if ok {
				a.nt.send(promise)
			}
		default:
			log.Panicf("Acceptor: %d unexpected Message type: ", a.id, m.Type)
		}
	}
}

// If an Acceptor receives a prepare request with number N greater
// than that of any prepare request To which it has already responded,
// then it responds To the request with a promise not To accept any more
// proposals numbered less than N and with the highest-numbered proposal
// (if any) that it has accepted.
func (a *Acceptor) receivePrepare(prepare Message) (Message, bool) {
	if a.promised.number() >= prepare.number() {
		log.Printf("Acceptor: %d [promised: %+v] ignored prepare %+v", a.id, a.promised, prepare)
		return Message{}, false
	}
	log.Printf("Acceptor: %d [promised: %+v] promised %+v", a.id, a.promised, prepare)
	a.promised = prepare
	m := Message{
		Type: Promise,
		From: a.id, To: prepare.From,
		N: a.promised.number(),
		// previously accepted proposal
		PrevN: a.accept.N, Value: a.accept.Value,
	}
	return m, true
}

// If an Acceptor receives an accept request for a proposal numbered
// N, it accepts the proposal unless it has already responded To a prepare
// request having a number greater than N.
func (a *Acceptor) receivePropose(propose Message) bool {
	if a.promised.number() > propose.number() {
		log.Printf("Acceptor: %d [promised: %+v] ignored proposal %+v", a.id, a.promised, propose)
		return false
	}
	if a.promised.number() < propose.number() {
		log.Panicf("Acceptor: %d received unexpected proposal %+v", a.id, propose)
	}
	log.Printf("Acceptor: %d [promised: %+v, accept: %+v] accepted proposal %+v", a.id, a.promised, a.accept, propose)
	a.accept = propose
	a.accept.Type = Accept
	return true
}

func (a *Acceptor) restart() {}
func (a *Acceptor) delay()   {}
