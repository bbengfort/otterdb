package events

const MaxAggregation = 512

// Runs a normal event loop handling one event at a time until the event channel is closed.
func Loop(events <-chan Event, handler Handler) (err error) {
	for e := range events {
		if err = handler.Handle(e); err != nil {
			return err
		}
	}
	return nil
}

// Runs an event loop that aggregates multiple write-ahead requests into a single append
// entries request to optimize distributed consensus and improve response times during
// high volume periods.
func AggregatingLoop(events <-chan Event, handler Handler) (err error) {
	for e := range events {
		if e.Event() == WriteAhead {
			// If we have a write-ahead request, attempt to aggregate, keeping track of
			// a next value (defaulting to nil) and storing all requests in an array
			// to be handled at once.
			var next Event
			requests := AggregatedWriteAheadEvents{e}

		aggregator:
			// The aggregator loop keeps reading events off the channel until there
			// is nothing on it or a non write-ahead event is read. In the meantime,
			// it aggregates all write-ahead events into a single events array.
			for {
				select {
				case next = <-events:
					if next.Event() != WriteAhead {
						break aggregator
					}
					requests = append(requests, next)
				default:
					// nothing is on the channel, break aggregator and do not handle
					// the empty next value by marking it as nil
					next = nil
					break aggregator
				}

				// If the requests array is too large, stop aggregating
				if len(requests) > MaxAggregation {
					next = nil
					break aggregator
				}
			}

			// This section happens after the aggregator for loop is complete.
			// First handle the write-ahead events, aggregating if there is more than
			// one write-ahead, otherwise handle normally.
			if len(requests) > 1 {
				if err = handler.Handle(requests); err != nil {
					return err
				}
			} else {
				// Handle the single write-ahead without the aggregator
				if err = handler.Handle(requests[0]); err != nil {
					return err
				}
			}

			// Second, handle the next event if one exists
			if next != nil {
				if err = handler.Handle(next); err != nil {
					return err
				}
			}

		} else {
			// Otherwise handle event normally without aggregation
			if err := handler.Handle(e); err != nil {
				return err
			}
		}
	}

	return nil
}
