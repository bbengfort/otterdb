package config

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
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
	Maintenance bool                `default:"false" desc:"if true, the node will start in maintenance mode"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info" desc:"specify the verbosity of logging (trace, debug, info, warn, error, fatal panic)"`
	ConsoleLog  bool                `split_words:"true" default:"false" desc:"if true logs colorized human readable output instead of json"`
	Server      ServerConfig
	Replica     ReplicaConfig
	Web         WebConfig
	processed   bool
}

type ServerConfig struct {
	Maintenance bool   `env:"OTTER_MAINTENANCE" desc:"if true sets the server to maintenance mode; inherited from parent"`
	Enabled     bool   `default:"true" desc:"if false, the client facing server will not be started, e.g. to uses this as a backup replica only"`
	BindAddr    string `default:":2202" split_words:"true" desc:"the ip address and port to bind the database server on"`
}

type ReplicaConfig struct {
	Maintenance bool   `env:"OTTER_MAINTENANCE" desc:"if true sets the replica to maintenance mode; inherited from parent"`
	Enabled     bool   `default:"false" desc:"if false, the replica service will not be started, e.g. run as a single node cluster"`
	BindAddr    string `default:":2204" split_words:"true" desc:"the ip address and port to bind the replica server on"`
	Aggregate   bool   `default:"true" desc:"if true the replica will aggregate append entries messages into a single consensus ballot"`
}

type WebConfig struct {
	Maintenance bool   `env:"OTTER_MAINTENANCE" desc:"if true sets the web ui to maintenance mode; inherited from parent"`
	Enabled     bool   `default:"false" desc:"set this to true to enable the web ui for the  database (opt-in)"`
	Mode        string `default:"release" desc:"specify the mode of the web server (release, debug, test)"`
	BindAddr    string `default:":2208" split_words:"true" desc:"the ip address and port to bind the web server on"`
	Origin      string `default:"http://localhost:2208" desc:"origin (url) of the web ui for creating endpoints and CORS access (include scheme, no trailing slash)"`
}

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
	if serr := c.Server.Validate(); serr != nil {
		err = errors.Join(err, serr)
	}

	if serr := c.Replica.Validate(); serr != nil {
		err = errors.Join(err, serr)
	}

	if serr := c.Web.Validate(); serr != nil {
		err = errors.Join(err, serr)
	}

	return err
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c ServerConfig) Validate() error {
	return nil
}

func (c ReplicaConfig) Validate() error {
	return nil
}

func (c WebConfig) Validate() (err error) {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		err = errors.Join(err, fmt.Errorf("invalid web configuration: %q is not a valid gin mode", c.Mode))
	}

	return err
}
