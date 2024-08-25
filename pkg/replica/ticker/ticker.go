package ticker

import (
	"context"
	"time"

	"github.com/bbengfort/otterdb/pkg/replica/events"
)

// Return a new ticker with the specified interval.
func New(interval Interval, event events.Event) *Ticker {
	return NewWithContext(context.Background(), interval, event)
}

// Return a new ticker with the specified context and interval. If this context is
// canceled, the ticker will automatically stop.
func NewWithContext(ctx context.Context, interval Interval, event events.Event) *Ticker {
	// Create an internal cancel function to stop the ticker
	ctx, cancel := context.WithCancel(ctx)

	// Create the internal send channel
	c := make(chan events.Event, 1)

	t := &Ticker{
		C:         c,
		interval:  interval,
		cancel:    cancel,
		interrupt: make(chan struct{}),
		event:     event,
	}

	go t.run(ctx, c)
	return t
}

// Ticker holds a channel that it uses to regularly dispatch time-based events on. The
// interval between events can be either fixed or stochastic (random) between delays.
// When the ticker is created it is started and it can only be stopped, at which point
// no new ticks will be sent, or it can be interrupted, which will reset the internal
// timer to perform a new tick after a new delay.
//
// The ticker will adjust the intervals or drop events to make up for slow receivers.
// TODO: send on an event channel rather than a time channel.
type Ticker struct {
	C         <-chan events.Event // The channel on which the events are delivered.
	interval  Interval
	cancel    context.CancelFunc
	interrupt chan struct{}
	event     events.Event
}

// Stop the ticker from sending any more events on the channel.
func (t *Ticker) Stop() {
	t.cancel()
}

// Stop the current interval and reset to the next interval of the ticker.
// Will panic if the ticker has been stopped.
func (t *Ticker) Interrupt() {
	t.interrupt <- struct{}{}
}

func (t *Ticker) Delay() time.Duration {
	return t.interval.Delay()
}

func (t *Ticker) run(ctx context.Context, c chan<- events.Event) {
	timer := time.NewTimer(t.interval.Delay())

	for {
		select {
		case <-timer.C:
			// Reset the internal timer for the next duration
			// Since the go routine has already received a value from timer.C
			// the timer is known to have expired and the channel drained, so
			// t.Reset can be used directly.
			timer.Reset(t.interval.Delay())

			// Non-blocking broadcast of event
			select {
			case c <- t.event:
			default:
			}

		case <-t.interrupt:
			// Stop the internal timer and drain its channel if needed
			if !timer.Stop() {
				<-timer.C
			}

			// Now that the timer has been drained, we can reset the timer
			timer.Reset(t.interval.Delay())

		case <-ctx.Done():
			// Stop the internal timer and drain its channel if needed
			if !timer.Stop() {
				<-timer.C
			}

			// Close the interrupt channel, draining it if necessary
			select {
			case <-t.interrupt:
			default:
			}

			// Close the channels to clean up resources
			close(t.interrupt)
			close(c)
			return
		}
	}
}
