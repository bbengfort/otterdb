package ticker

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	mrand "math/rand"
)

// Random source used to generate pseudo-random numbers for the ticker.
var rand *mrand.Rand

// Max attempts for non-negative random number generation.
var maxAttempts = 4

func init() {
	ResetSource()
}

// SetSource allows the user to specify a new pseudo-random source for random number
// generation in this package; e.g. for deterministic unit testing.
func SetSource(s mrand.Source) {
	rand = mrand.New(s)
}

// ResetSource re-initializes the pseudo-random source; seeding with a random value.
func ResetSource() {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		panic("cryptographically random number generator required to seed source")
	}
	rand = mrand.New(mrand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}

// Returns the min and max values for n after applying scaling factor s.
// Protects from max overflow by truncating the upper bound and returning max int64.
func scaleBounds(n int64, f float64) (min, max int64) {
	minf := math.Floor(float64(n) * (1 - f))
	maxf := math.Ceil(float64(n) * (1 + f))

	if maxf > math.MaxInt64 {
		return int64(minf), math.MaxInt64
	}

	return int64(minf), int64(maxf)
}

// Returns a non-negative pseudo-random number in the half open interval [min, max).
func randRange(min, max int64) int64 {
	if min == max {
		return min
	}
	return rand.Int63n(max-min) + min
}

// Returns a non-negative pseudo-random number from a normal distribution specified by
// the mean and standard deviation. The function attempts multiple times to generate a
// non-negative random number, then simply returns the mean. This means that a true
// normal distribution is not returned, particularly if the mean is close to zero with
// respect to the standard deviation.
func randNormal(mean int64, sdev float64) int64 {
	for i := 0; i < maxAttempts; i++ {
		sample := rand.NormFloat64()*sdev + float64(mean)
		if sample > 0.0 {
			return int64(math.Round(sample))
		}
	}

	// After max attempts, simply return the mean.
	return mean
}
