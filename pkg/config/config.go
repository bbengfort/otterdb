package config

import (
	"github.com/rotationalio/confire"
	"github.com/rs/zerolog"

	"github.com/bbengfort/otterdb/pkg/logger"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be OTTER_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const Prefix = "otter"

// Config contains all of the configuration parameters for an otterdb instance and is
// loaded from the environment or a configuration file with reasonable defaults for
// values that are omitted. The Config should be validated in preparation for running
// the otterdb instance to ensure that all server operations work as expected.
type Config struct {
	Maintenance bool                `default:"false" yaml:"maintenance"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info" yaml:"log_level"`
	ConsoleLog  bool                `split_words:"true" default:"false" yaml:"console_log"`
	Server      ServerConfig
	Replica     ReplicaConfig
	Web         WebConfig
	processed   bool
}

type ServerConfig struct{}

type ReplicaConfig struct{}

type WebConfig struct{}

func New() (conf Config, err error) {
	if err = confire.Process(Prefix, &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Returns true if the config has not been correctly processed from the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Custom validations are added here, particularly validations that require one or more
// fields to be processed before the validation occurs.
// NOTE: ensure that all nested config validation methods are called here.
func (c Config) Validate() (err error) {
	// if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
	// 	return fmt.Errorf("invalid configuration: %q is not a valid gin mode", c.Mode)
	// }
	return nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}
