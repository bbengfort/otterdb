package events

// Events can be used to pass datagrams to the one-big pipe channel for distributed
// consensus processing. Some events are just empty types, but others may contain data
// such as votes, entries to append, or indices to commit.
type Event interface {
	Event() EventType
}

// Handler is a function that takes an event and does something with it,
// possibly returning an error. If an event has multiple handlers or callbacks,
// the semantics are that they are each called synchronously in order; if an
// error is returned from one of the handlers, then none of the remaining
// handlers are called.
type Handler interface {
	Handle(e Event) error
}

// BufferSize for a "one big pipe" event channel.
const BufferSize = 1024

// EventType allows callers to detect the kind of event without a type assertion.
type EventType uint8

const (
	UnknownEvent EventType = iota
	ErrorEvent
	WriteAhead
	AggregatedWriteAhead
	HeartbeatTimeout
	ElectionTimeout
	VoteRequest
	VoteReply
	AppendRequest
	AppendReply
)

// Names of event types for easier debugging
var eventTypes = [...]string{
	"unknown", "error",
	"writeAhead", "aggregatedWriteAhead",
	"heartbeatTimeout", "electionTimeout",
	"voteRequest", "voteReply",
	"appendRequest", "appendReply",
}

func (t EventType) String() string {
	if int(t) < len(eventTypes) {
		return eventTypes[t]
	}
	return eventTypes[0]
}

// TODO: replace with actual data structure
type AggregatedWriteAheadEvents []Event

func (a AggregatedWriteAheadEvents) Event() EventType {
	return AggregatedWriteAhead
}
