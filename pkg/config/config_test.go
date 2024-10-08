package config_test

import (
	"os"
	"testing"

	"github.com/bbengfort/otterdb/pkg/config"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"OTTER_MAINTENANCE":       "true",
	"OTTER_LOG_LEVEL":         "debug",
	"OTTER_CONSOLE_LOG":       "true",
	"OTTER_SERVER_ENABLED":    "false",
	"OTTER_SERVER_BIND_ADDR":  ":3303",
	"OTTER_REPLICA_ENABLED":   "true",
	"OTTER_REPLICA_BIND_ADDR": ":3304",
	"OTTER_REPLICA_AGGREGATE": "false",
	"OTTER_WEB_ENABLED":       "true",
	"OTTER_WEB_MODE":          "test",
	"OTTER_WEB_BIND_ADDR":     ":3305",
	"OTTER_WEB_ORIGIN":        "https://example.com",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after the test is complete.
	t.Cleanup(cleanupEnv())
	setEnv()

	conf, err := config.New()
	require.NoError(t, err, "could not process configuration from the environment")
	require.False(t, conf.IsZero(), "processed config should not be zero valued")

	// Ensure configuration is correctly set from the environment
	require.True(t, conf.Maintenance)
	require.True(t, conf.Server.Maintenance)
	require.True(t, conf.Replica.Maintenance)
	require.True(t, conf.Web.Maintenance)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.False(t, conf.Server.Enabled)
	require.Equal(t, testEnv["OTTER_SERVER_BIND_ADDR"], conf.Server.BindAddr)
	require.True(t, conf.Replica.Enabled)
	require.Equal(t, testEnv["OTTER_REPLICA_BIND_ADDR"], conf.Replica.BindAddr)
	require.False(t, conf.Replica.Aggregate)
	require.True(t, conf.Web.Enabled)
	require.Equal(t, testEnv["OTTER_WEB_MODE"], conf.Web.Mode)
	require.Equal(t, testEnv["OTTER_WEB_BIND_ADDR"], conf.Web.BindAddr)
	require.Equal(t, testEnv["OTTER_WEB_ORIGIN"], conf.Web.Origin)
}

// Returns the current environment for the specified keys, or if no keys are specified
// then it returns the current environment for all keys in the testEnv variable.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv variable. If no keys are specified,
// then this function sets all environment variables from the testEnv.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}

// Cleanup helper function that can be run when the tests are complete to reset the
// environment back to its previous state before the test was run.
func cleanupEnv(keys ...string) func() {
	prevEnv := curEnv(keys...)
	return func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
