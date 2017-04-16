package paxos

type msgType int

const (
	Prepare msgType = iota + 1
	Propose
	Promise
	Accept
)

type Message struct {
	From, To int
	Type     msgType
	N        int
	PrevN    int
	Value    string
}

func (m Message) number() int {
	return m.N
}

func (m Message) proposalValue() string {
	switch m.Type {
	case Promise, Accept:
		return m.Value
	default:
		panic("unexpected proposalV")
	}
}

func (m Message) proposalNumber() int {
	switch m.Type {
	case Promise:
		return m.PrevN
	case Accept:
		return m.N
	default:
		panic("unexpected proposalN")
	}
}

type promise interface {
	number() int
}

type accept interface {
	proposalValue() string
	proposalNumber() int
}
