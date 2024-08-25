package replica

import "github.com/bbengfort/otterdb/pkg/replica/events"

func (r *Replica) Handle(e events.Event) error {
	return nil
}

func (r *Replica) Dispatch(e events.Event) error {
	if r.events == nil {
		return ErrNotListening
	}

	r.events <- e
	return nil
}
