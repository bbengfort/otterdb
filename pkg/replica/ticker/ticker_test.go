package ticker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/bbengfort/otterdb/pkg/replica/ticker"
	"github.com/stretchr/testify/require"
)

func TestTicker(t *testing.T) {
	var wg sync.WaitGroup
	beat := ticker.New(ticker.Fixed(50*time.Millisecond), ticker.HeartbeatTimeout{})

	// Create counter routine
	wg.Add(1)
	ticks := 0

	go func() {
		defer wg.Done()
		for range beat.C {
			ticks++
		}
	}()

	// Create stopping routine
	wg.Add(1)
	time.AfterFunc(275*time.Millisecond, func() {
		beat.Stop()
		wg.Done()
	})

	// Ensure that the delay function connects to the interval
	require.Equal(t, 50*time.Millisecond, beat.Delay())

	wg.Wait()
	require.Equal(t, 5, ticks)
}

func TestInterrupt(t *testing.T) {
	var wg sync.WaitGroup
	beat := ticker.New(ticker.Fixed(50*time.Millisecond), ticker.ElectionTimeout{})

	// Create counter routine
	wg.Add(1)
	ticks := 0

	go func() {
		defer wg.Done()
		for range beat.C {
			ticks++
		}
	}()

	// Create interrupting routines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		delay := time.Duration(int64(i+1)*int64(45)) * time.Millisecond
		time.AfterFunc(delay, func() {
			beat.Interrupt()
			wg.Done()
		})
	}

	// Create stopping routine
	wg.Add(1)
	time.AfterFunc(255*time.Millisecond, func() {
		beat.Stop()
		wg.Done()
	})

	wg.Wait()
	require.Equal(t, 0, ticks)
}
