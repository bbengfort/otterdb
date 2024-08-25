package events_test

import (
	"fmt"
	"testing"

	"github.com/bbengfort/otterdb/pkg/replica/events"
	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) {
	testCases := []struct {
		event    fmt.Stringer
		expected string
	}{
		{events.UnknownEvent, "unknown"},
		{events.ErrorEvent, "error"},
		{events.WriteAhead, "writeAhead"},
		{events.HeartbeatTimeout, "heartbeatTimeout"},
		{events.ElectionTimeout, "electionTimeout"},
		{events.VoteRequest, "voteRequest"},
		{events.VoteReply, "voteReply"},
		{events.AppendRequest, "appendRequest"},
		{events.AppendReply, "appendReply"},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, tc.event.String(), "test case %d failed", i)
	}
}
