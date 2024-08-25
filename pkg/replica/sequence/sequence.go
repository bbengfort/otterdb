package sequence

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
)

// Sequence implements a thread-safe monotonically increasing counter. This simple data
// structure does not take into account integer overflow, allowing the counter to roll
// over back to zero to start the sequence again. The sequence can also be saved to
// disk as binary data, text data, or json data.
type Sequence struct {
	sync.Mutex
	counter uint64
}

// Create a new sequence
func New() *Sequence {
	return &Sequence{}
}

// Create a new sequence starting at the specified value.
func Start(at uint64) *Sequence {
	return &Sequence{counter: at}
}

func (s *Sequence) Next() uint64 {
	s.Lock()
	defer s.Unlock()
	s.counter++
	return s.counter
}

func (s *Sequence) MarshalBinary() ([]byte, error) {
	data := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(data, s.counter)
	return data[:n], nil
}

func (s *Sequence) UnmarshalBinary(data []byte) error {
	var n int
	s.counter, n = binary.Uvarint(data)
	if n != len(data) {
		return errors.New("sequence overflowed")
	}
	return nil
}

func (s *Sequence) MarshalText() ([]byte, error) {
	text := strconv.FormatUint(s.counter, 16)
	return []byte(text), nil
}

func (s *Sequence) UnmarshalText(text []byte) (err error) {
	if s.counter, err = strconv.ParseUint(string(text), 16, 64); err != nil {
		return err
	}
	return nil
}

func (s *Sequence) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.counter)
}

func (s *Sequence) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.counter)
}
