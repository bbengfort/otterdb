package ticker

import (
	"errors"
	"time"
)

// An interval specifies the delay between ticks for a ticker. The delay can be fixed so
// that the same delay is returned every time, or it can be stochastic, meaning that a
// random delay is computed for every tick.
type Interval interface {
	Delay() time.Duration
}

// Returns a fixed interval that always returns the specified delay.
func Fixed(delay time.Duration) Interval {
	return fixed{delay: delay}
}

type fixed struct {
	delay time.Duration
}

func (f fixed) Delay() time.Duration {
	return f.delay
}

// Returns a uniform random delay in the half open interval [min, max) delay.
func Uniform(min, max time.Duration) Interval {
	return uniform{int64(min), int64(max)}
}

// Jitter scales the specified duration within a factor f and returns a random delay
// around that factor. E.g. if your delay is 5s with a factor of .25 then the delay
// will be random in the range 3.75s to 6.25s. The smaller the factor the smaller the
// range of the jitter, the larger the factor, the larger the range of the jitter.
//
// Duration must be greater than zero and the scaling factor f must be within the range
// 0 < f <= 1.0, otherwise this function will panic.
func Jitter(delay time.Duration, f float64) Interval {
	switch {
	case delay <= 0:
		panic(errors.New("non-positive or zero interval for delay"))
	case f > 1.0 || f <= 0:
		panic(errors.New("scaling factor must be 0 < f <= 1.0"))
	}

	min, max := scaleBounds(int64(delay), f)
	return uniform{min, max}
}

type uniform struct {
	min int64
	max int64
}

func (u uniform) Delay() time.Duration {
	return time.Duration(randRange(u.min, u.max))
}

func Normal(mean, stddev time.Duration) Interval {
	switch {
	case mean <= 0:
		panic(errors.New("non-positive or zero interval for delay"))
	case stddev <= 0:
		panic(errors.New("non-positive or zero value for standard deviation"))
	}

	return normal{int64(mean), float64(stddev)}
}

type normal struct {
	mean int64
	sdev float64
}

func (n normal) Delay() time.Duration {
	return time.Duration(randNormal(n.mean, n.sdev))
}
