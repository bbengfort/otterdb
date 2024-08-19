package logger_test

import (
	"context"
	"testing"

	"github.com/bbengfort/otterdb/pkg/logger"

	"github.com/stretchr/testify/require"
)

func TestRequestIDContext(t *testing.T) {
	requestID := "01J5N8FV669X7WE1FY7SSPEY1T"
	parent, cancel := context.WithCancel(context.Background())
	ctx := logger.WithRequestID(parent, requestID)

	cmp, ok := logger.RequestID(ctx)
	require.True(t, ok)
	require.Equal(t, requestID, cmp)

	cancel()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}
