package quorum_test

import (
	"fmt"
	"math"
	"testing"

	. "github.com/bbengfort/otterdb/pkg/replica/quorum"
	"github.com/stretchr/testify/require"
)

func TestElection(t *testing.T) {
	t.Run("Panics", func(t *testing.T) {
		quorum := New()
		require.PanicsWithError(t, ErrNoEmptyQuorum.Error(), func() { quorum.Election() })

		hosts := make([]string, 0, math.MaxUint16+1)
		for i := 0; i < math.MaxUint16+1; i++ {
			hosts = append(hosts, fmt.Sprintf("q%x", i))
		}
		quorum = New(hosts...)
		require.PanicsWithError(t, ErrQuorumTooLarge.Error(), func() { quorum.Election() })
	})

	t.Run("VoteValidation", func(t *testing.T) {
		quorum := New("jade", "kira", "opal")
		election := quorum.Election()

		_, err := election.Vote("artemis")
		require.EqualError(t, err, `"artemis" is not a member of the quorum`)

		election.Vote("jade")
		_, err = election.Vote("jade")
		require.EqualError(t, err, `"jade" has already voted in this election`)
	})

	t.Run("NoCheating", func(t *testing.T) {
		quorum := New("jade", "kira", "opal")
		election := quorum.Election()

		passed, err := election.Vote("kira")
		require.False(t, passed)
		require.NoError(t, err)

		passed, err = election.Vote("artemis")
		require.False(t, passed)
		require.Error(t, err)

		passed, err = election.Vote("kira")
		require.False(t, passed)
		require.Error(t, err)

		passed, err = election.Vote("opal")
		require.True(t, passed)
		require.NoError(t, err)

		passed, err = election.Vote("opal")
		require.False(t, passed)
		require.Error(t, err)
	})
}

func TestElectionSize(t *testing.T) {
	// Make a set of host ids to reuse accross tests.
	hosts := make([]string, 256)
	for i := range hosts {
		hosts[i] = fmt.Sprintf("seat%03d", i+1)
	}

	makeTest := func(size, majority int) func(t *testing.T) {
		return func(t *testing.T) {
			quorum := New(hosts[:size]...)
			election := quorum.Election()

			// Vote should not pass up to the majority
			for i := 0; i < majority-1; i++ {
				passed, err := election.Vote(hosts[i])
				require.NoError(t, err)
				require.False(t, passed)
			}

			// Vote should pass when the majority is reached
			passed, err := election.Vote(hosts[majority])
			require.NoError(t, err)
			require.True(t, passed)

			// Vote should remain passed after the majority is reached
			for i := majority + 1; i < size; i++ {
				passed, err := election.Vote(hosts[i])
				require.NoError(t, err)
				require.True(t, passed)
			}

		}
	}

	t.Run("Q3", makeTest(3, 2))
	t.Run("Q4", makeTest(4, 3))
	t.Run("Q5", makeTest(5, 3))
	t.Run("Q6", makeTest(6, 4))
	t.Run("Q7", makeTest(7, 4))
	t.Run("Q9", makeTest(9, 5))
	t.Run("Q11", makeTest(11, 6))
	t.Run("Q13", makeTest(13, 7))
	t.Run("Q99", makeTest(99, 50))
	t.Run("Q256", makeTest(256, 129))
}
