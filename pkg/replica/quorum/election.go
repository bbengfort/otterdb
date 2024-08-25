package quorum

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrNoEmptyQuorum  = errors.New("cannot create a vote for an empty quorum")
	ErrQuorumTooLarge = errors.New("cannot create a vote for extremely large quorum")
)

// NewElection creates and initializes a vote data structure for the specified
// quorum. A vote is a throw away counter of the number of accept votes received,
// ensuring that only members of the quorum are allowed to vote.
func NewElection(quorum *Quorum) *Election {
	size := len(quorum.hosts)
	switch {
	case size < 1:
		panic(ErrNoEmptyQuorum)
	case size > math.MaxUint16:
		panic(ErrQuorumTooLarge)
	}

	// Track the hosts that have already cast votes
	q := make(map[string]bool, len(quorum.hosts))
	for host := range quorum.hosts {
		q[host] = false
	}

	return &Election{quorum: q}
}

// Election implements a decision making data structure for a specific quorum
// such that each host can cast an accept vote. Once a majority of the quorum casts
// accept ballots the vote is considered passed.
type Election struct {
	quorum  map[string]bool
	ballots uint16
}

// Vote records the vote for the given member, identified by name; will
// return an error if the voting member is not part of the quorum or if the member has
// already cast a vote in the specified election. Passed will return true if the number
// of ballots is greater than or equal to the majority (n/2+1).
func (e *Election) Vote(member string) (passed bool, err error) {
	if vote, isMember := e.quorum[member]; !isMember {
		return false, fmt.Errorf("%q is not a member of the quorum", member)
	} else if vote {
		return false, fmt.Errorf("%q has already voted in this election", member)
	}

	e.quorum[member] = true
	e.ballots++

	return e.ballots >= uint16((len(e.quorum)/2)+1), nil
}
