package peers_test

import (
	"path/filepath"
	"testing"

	. "github.com/bbengfort/otterdb/pkg/replica/peers"
	"github.com/stretchr/testify/require"
)

func TestPeers(t *testing.T) {
	peers, err := Load("testdata/peers.json")
	require.NoError(t, err, "could not load testdata peers")

	t.Run("Names", func(t *testing.T) {
		names := peers.Names()
		require.ElementsMatch(t, names, []string{"opal", "jade", "kira"})
	})

	t.Run("Get", func(t *testing.T) {
		peer, err := peers.Get("jade")
		require.NoError(t, err, "could not get peer")
		require.Equal(t, uint16(30), peer.PID)

		peer, err = peers.Get("artemis")
		require.EqualError(t, err, "no peer found named \"artemis\"")
		require.Nil(t, peer)
	})

	t.Run("Presiding", func(t *testing.T) {
		require.Equal(t, "kira", peers.Presiding())
	})
}

func TestSerialization(t *testing.T) {
	peers := Peers{
		{
			PID:    20,
			Name:   "opal",
			Addr:   "opal.io:443",
			Region: "us-central1",
		},
		{
			PID:    30,
			Name:   "kira",
			Addr:   "kira.com:443",
			Region: "eu-west1",
		},
		{
			PID:    40,
			Name:   "jade",
			Addr:   "jade.co:443",
			Region: "us-east4",
		},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "testpeers.json")

	err := peers.Dump(path)
	require.NoError(t, err, "could not dump peers")

	cmp, err := Load(path)
	require.NoError(t, err, "could not load peers")

	require.Equal(t, peers, cmp)
}
