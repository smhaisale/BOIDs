package paxos

import (
	"log"
	"time"
)

type Proposer struct {
	id int
	// stable
	lastSeq int

	value  string
	valueN int

	acceptors map[int]promise
	nt        Network
}

func newProposer(id int, value string, nt Network, acceptors ...int) *Proposer {
	p := &Proposer{id: id, nt: nt, lastSeq: 0, value: value, acceptors: make(map[int]promise)}
	for _, a := range acceptors {
		p.acceptors[a] = Message{}
	}
	return p
}

func (p *Proposer) run() {
	var ok bool
	var m Message

	// stage 1: do prepare until reach the majority
	for !p.majorityReached() {
		if !ok {
			ms := p.prepare()
			for i := range ms {
				p.nt.send(ms[i])
			}
		}
		m, ok = p.nt.recv(time.Second)
		if !ok {
			// the previous prepare is failed
			// continue To do another prepare
			continue
		}

		switch m.Type {
		case Promise:
			p.receivePromise(m)
		default:
			log.Panicf("Proposer: %d unexpected Message type: ", p.id, m.Type)
		}
	}
	log.Printf("Proposer: %d promise %d reached majority %d", p.id, p.n(), p.majority())

	// stage 2: do propose
	log.Printf("Proposer: %d starts To propose [%d: %s]", p.id, p.n(), p.value)
	ms := p.propose()
	for i := range ms {
		p.nt.send(ms[i])
	}
}

// If the Proposer receives the requested responses From a majority of
// the acceptors, then it can issue a proposal with number N and Value
// v, where v is the Value of the highest-numbered proposal among the
// responses, or is any Value selected by the Proposer if the responders
// reported no proposals.
func (p *Proposer) propose() []Message {
	ms := make([]Message, p.majority())

	i := 0
	for to, promise := range p.acceptors {
		if promise.number() == p.n() {
			ms[i] = Message{From: p.id, To: to, Type: Propose, N: p.n(), Value: p.value}
			i++
		}
		if i == p.majority() {
			break
		}
	}
	return ms
}

// A Proposer chooses a new proposal number N and sends a request To
// each member of some set of acceptors, asking it To respond with:
// (a) A promise never again To accept a proposal numbered less than N, and
// (b) The proposal with the highest number less than N that it has accepted, if any.
func (p *Proposer) prepare() []Message {
	p.lastSeq++

	ms := make([]Message, p.majority())
	i := 0
	for to := range p.acceptors {
		ms[i] = Message{From: p.id, To: to, Type: Prepare, N: p.n()}
		i++
		if i == p.majority() {
			break
		}
	}
	return ms
}

func (p *Proposer) receivePromise(promise Message) {
	prevPromise := p.acceptors[promise.From]

	if prevPromise.number() < promise.number() {
		log.Printf("Proposer: %d received a new promise %+v", p.id, promise)
		p.acceptors[promise.From] = promise

		//update Value To the Value with a larger N
		if promise.proposalNumber() > p.valueN {
			log.Printf("Proposer: %d updated the Value [%s] To %s", p.id, p.value, promise.proposalValue())
			p.valueN = promise.proposalNumber()
			p.value = promise.proposalValue()
		}
	}
}

func (p *Proposer) majority() int { return len(p.acceptors)/2 + 1 }

func (p *Proposer) majorityReached() bool {
	m := 0
	for _, promise := range p.acceptors {
		if promise.number() == p.n() {
			m++
		}
	}
	if m >= p.majority() {
		return true
	}
	return false
}

func (p *Proposer) n() int { return p.lastSeq<<16 | p.id }
