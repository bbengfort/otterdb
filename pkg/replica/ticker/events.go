package ticker

import "github.com/bbengfort/otterdb/pkg/replica/events"

type HeartbeatTimeout struct{}

var heartbeatTick = HeartbeatTimeout{}

func (t HeartbeatTimeout) Event() events.EventType {
	return events.HeartbeatTimeout
}

type ElectionTimeout struct{}

var electionTimeout = ElectionTimeout{}

func (t ElectionTimeout) Event() events.EventType {
	return events.ElectionTimeout
}

func NewHeartbeatTicker(interval Interval) *Ticker {
	return New(interval, heartbeatTick)
}

func NewElectionTicker(interval Interval) *Ticker {
	return New(interval, electionTimeout)
}
