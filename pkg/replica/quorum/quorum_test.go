package quorum_test

import (
	"testing"

	. "github.com/bbengfort/otterdb/pkg/replica/quorum"
	"github.com/stretchr/testify/require"
)

func TestQuorum(t *testing.T) {

	t.Run("Create", func(t *testing.T) {
		hosts := []string{"192.168.1.2:3265", "192.168.1.3:3264", "192.168.1.2:3264"}
		quorum := New(hosts...)

		require.Greater(t, quorum.ID(), uint64(0))
		require.Equal(t, 3, quorum.Size())
		require.ElementsMatch(t, hosts, quorum.Hosts())
	})

	t.Run("Unique", func(t *testing.T) {
		quorum := New("192.168.1.1", "192.168.1.2", "192.168.1.2", "192.168.1.1")
		require.Equal(t, 2, quorum.Size())
		require.ElementsMatch(t, []string{"192.168.1.1", "192.168.1.2"}, quorum.Hosts())
	})

	t.Run("Independent", func(t *testing.T) {
		q1 := New("192.168.1.1", "192.168.1.2", "192.168.1.3")
		q2 := New("192.168.1.1", "192.168.1.4", "192.168.1.5")
		q3 := New("192.168.1.4", "192.168.1.2", "192.168.1.1")
		q4 := New("192.168.1.1", "192.168.1.2", "192.168.1.3")

		require.NotEqual(t, q1.ID(), q2.ID())
		require.NotEqual(t, q2.ID(), q3.ID())
		require.NotEqual(t, q3.ID(), q4.ID())
	})

	t.Run("Contains", func(t *testing.T) {
		q1 := New("192.168.1.1", "192.168.1.2", "192.168.1.3")
		require.True(t, q1.Contains("192.168.1.1"))
		require.True(t, q1.Contains("192.168.1.2"))
		require.True(t, q1.Contains("192.168.1.3"))
		require.False(t, q1.Contains("192.168.1.4"))
		require.False(t, q1.Contains(""))
	})

	t.Run("NoEmptyHost", func(t *testing.T) {
		require.False(t, New("").Contains(""))
		require.False(t, New().Contains(""))
	})

	t.Run("Subset", func(t *testing.T) {
		// master set
		q1 := New("a", "b", "c", "d", "e", "f", "g", "h")

		// is subset test cases
		for _, q := range []*Quorum{
			New("a", "b", "c"),
			New("f"),
			New("c", "d", "b", "g", "h"),
		} {
			require.True(t, q.IsSubset(q1), "correct subset failed")
		}

		require.True(t, q1.IsSubset(q1), "master should be subset of self")
		require.True(t, New().IsSubset(q1), "empty should be a subset of quorum")

		// not subset test cases
		for _, q := range []*Quorum{
			New("a", "b", "c", "z"),
			New("s", "t", "r"),
			New("c", "q", "B", "O", "z"),
		} {
			require.False(t, q.IsSubset(q1), "incorrect subset failed")
		}
	})

	t.Run("Superset", func(t *testing.T) {
		// master set
		q1 := New("a", "b", "c", "d", "e", "f", "g", "h")

		//is superset test cases
		for _, q := range []*Quorum{
			New("a", "b", "c"),
			New("f"),
			New("c", "d", "b", "g", "h"),
		} {
			require.True(t, q1.IsSuperset(q), "correct superset failed")
		}

		require.True(t, q1.IsSuperset(q1), "quorum is not superset of self")
		require.False(t, New().IsSuperset(q1), "empty is superset of quorum")
		require.True(t, q1.IsSuperset(New()), "quorum is not superset of empty")

		// not superset test cases
		for _, q := range []*Quorum{
			New("a", "b", "c", "z"),
			New("s", "t", "r"),
			New("c", "q", "B", "O", "z"),
		} {
			require.False(t, q1.IsSuperset(q), "incorrect superset failed")
		}
	})

	t.Run("Intersect", func(t *testing.T) {
		q1 := New("a", "b", "c")
		q2 := New("d", "e", "f")
		q3 := New("c", "e")
		empty := New()

		require.False(t, q1.Intersects(q2))
		require.True(t, q1.Intersects(q3))
		require.False(t, q1.Intersects(empty))
		require.False(t, q2.Intersects(q1))
		require.True(t, q2.Intersects(q3))
		require.False(t, q2.Intersects(empty))
		require.True(t, q3.Intersects(q1))
		require.True(t, q3.Intersects(q2))
		require.False(t, q3.Intersects(empty))
		require.False(t, empty.Intersects(q1))
		require.False(t, empty.Intersects(q2))
		require.False(t, empty.Intersects(q3))
	})
}
