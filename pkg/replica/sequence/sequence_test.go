package sequence_test

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"testing"

	"github.com/bbengfort/otterdb/pkg/replica/sequence"
	"github.com/stretchr/testify/require"
)

func TestSequence(t *testing.T) {
	var wg sync.WaitGroup
	seq := sequence.New()

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 256; i++ {
				seq.Next()
			}
		}()
	}

	wg.Wait()
	require.Equal(t, uint64(4097), seq.Next())
}

func TestSerialization(t *testing.T) {
	t.Run("Binary", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			seq := sequence.Start(randUint64())
			data, err := seq.MarshalBinary()
			require.NoError(t, err, "could not marshal binary")

			cmp := &sequence.Sequence{}
			err = cmp.UnmarshalBinary(data)
			require.NoError(t, err, "could not unmarshal binary")

			require.Equal(t, seq.Next(), cmp.Next())
		}
	})

	t.Run("Text", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			seq := sequence.Start(randUint64())
			data, err := seq.MarshalText()
			require.NoError(t, err, "could not marshal text")

			cmp := &sequence.Sequence{}
			err = cmp.UnmarshalText(data)
			require.NoError(t, err, "could not unmarshal text")

			require.Equal(t, seq.Next(), cmp.Next())
		}
	})

	t.Run("JSON", func(t *testing.T) {
		for i := 0; i < 256; i++ {
			seq := sequence.Start(randUint64())
			data, err := seq.MarshalJSON()
			require.NoError(t, err, "could not marshal json")

			cmp := &sequence.Sequence{}
			err = cmp.UnmarshalJSON(data)
			require.NoError(t, err, "could not unmarshal json")

			require.Equal(t, seq.Next(), cmp.Next())
		}
	})
}

func randUint64() uint64 {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(buf)
}
