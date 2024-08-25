package ticker_test

import (
	"testing"

	"github.com/bbengfort/otterdb/pkg/replica/events"
	"github.com/bbengfort/otterdb/pkg/replica/ticker"
	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) {
	t.Run("Heartbeat", func(t *testing.T) {
		require.Equal(t, ticker.HeartbeatTimeout{}.Event(), events.HeartbeatTimeout)
	})

	t.Run("Election", func(t *testing.T) {
		require.Equal(t, ticker.ElectionTimeout{}.Event(), events.ElectionTimeout)
	})
}
