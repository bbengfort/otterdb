package ticker_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/bbengfort/otterdb/pkg/replica/ticker"
	"github.com/stretchr/testify/require"
)

var maxTests = 256

func TestInterval(t *testing.T) {
	t.Run("Fixed", func(t *testing.T) {
		delay := 300 * time.Millisecond
		interval := ticker.Fixed(delay)
		for i := 0; i < maxTests; i++ {
			require.Equal(t, delay, interval.Delay())
		}
	})

	t.Run("Uniform", func(t *testing.T) {
		min := 300 * time.Millisecond
		max := 700 * time.Millisecond
		interval := ticker.Uniform(min, max)

		prev := time.Duration(0)
		for i := 0; i < maxTests; i++ {
			delay := interval.Delay()
			require.GreaterOrEqual(t, delay, min)
			require.Less(t, delay, max)

			require.NotEqual(t, prev, delay)
			prev = delay
		}
	})

	t.Run("Jitter", func(t *testing.T) {
		min := 3750 * time.Millisecond
		max := 6250 * time.Millisecond
		interval := ticker.Jitter(5*time.Second, 0.25)

		prev := time.Duration(0)
		for i := 0; i < maxTests; i++ {
			delay := interval.Delay()
			require.GreaterOrEqual(t, delay, min)
			require.Less(t, delay, max)

			require.NotEqual(t, prev, delay)
			prev = delay
		}
	})

	t.Run("Normal", func(t *testing.T) {
		min := 4800 * time.Millisecond
		max := 5200 * time.Millisecond
		interval := ticker.Normal(5*time.Second, 30*time.Millisecond)

		prev := time.Duration(0)
		for i := 0; i < maxTests; i++ {
			delay := interval.Delay()
			require.GreaterOrEqual(t, delay, min)
			require.Less(t, delay, max)

			require.NotEqual(t, prev, delay)
			prev = delay
		}
	})
}

func TestFixedInterval(t *testing.T) {
	defer ticker.ResetSource()

	makeTest := func(interval ticker.Interval, values []time.Duration) func(t *testing.T) {
		ticker.SetSource(rand.NewSource(42))
		return func(t *testing.T) {
			for _, val := range values {
				require.Equal(t, val, interval.Delay())
			}
		}
	}

	t.Run("Fixed", makeTest(ticker.Fixed(300*time.Millisecond), []time.Duration{300000000, 300000000, 300000000, 300000000, 300000000, 300000000, 300000000, 300000000, 300000000}))
	t.Run("Uniform", makeTest(ticker.Uniform(300*time.Microsecond, 800*time.Microsecond), []time.Duration{578675, 656411, 678760, 424009, 347657, 537261, 313247, 414208, 792868}))
	t.Run("Jitter", makeTest(ticker.Jitter(300*time.Microsecond, 0.25), []time.Duration{253675, 331411, 353760, 299009, 272657, 312261, 338247, 239208, 367868}))
	t.Run("Normal", makeTest(ticker.Normal(3000*time.Microsecond, 40*time.Microsecond), []time.Duration{3062145, 3005010, 2980225, 3049761, 3005279, 3048255, 2974970, 3025184, 3062621}))
}

func TestJitterPanic(t *testing.T) {
	require.PanicsWithError(t, "non-positive or zero interval for delay", func() { ticker.Jitter(0, 0.5) })
	require.PanicsWithError(t, "scaling factor must be 0 < f <= 1.0", func() { ticker.Jitter(10000, -0.5) })
}

func TestNormalPanic(t *testing.T) {
	require.PanicsWithError(t, "non-positive or zero interval for delay", func() { ticker.Normal(-100, -4) })
	require.PanicsWithError(t, "non-positive or zero value for standard deviation", func() { ticker.Normal(10000, 0.000) })
}
