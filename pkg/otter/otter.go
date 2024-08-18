/*
Package otter implements the otterdb service and is the primary entrypoint to the code.
*/
package otter

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bbengfort/otterdb/pkg/config"
	"github.com/bbengfort/otterdb/pkg/logger"
	"github.com/bbengfort/otterdb/pkg/replica"
	"github.com/bbengfort/otterdb/pkg/server"
	"github.com/bbengfort/otterdb/pkg/web"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Initializes zerolog with our default logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

type OtterDB struct {
	conf    config.Config
	replica *replica.Replica
	server  *server.Server
	web     *web.Server
	errc    chan error
}

func New(conf config.Config) (svc *OtterDB, err error) {
	// Load the default configuration from the environment if config is empty.
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Setup our logging config first thing
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Create the otterdb service
	svc = &OtterDB{conf: conf, errc: make(chan error, 1)}

	// Configure the replica service
	if svc.replica, err = replica.New(conf.Replica); err != nil {
		return nil, err
	}

	// Configure the database service
	if svc.server, err = server.New(conf.Server); err != nil {
		return nil, err
	}

	// Configure the web user interface service
	if svc.web, err = web.New(conf.Web); err != nil {
		return nil, err
	}

	return svc, nil
}

// Serve all enabled services based on configuration and block until shutdown.
func (o *OtterDB) Serve() (err error) {
	// Handle OS Signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		o.errc <- o.Shutdown()
	}()

	// Start the replica server first
	if err = o.replica.Serve(o.errc); err != nil {
		return err
	}

	// Start the database server for clients
	if err = o.server.Serve(o.errc); err != nil {
		return err
	}

	// Start the web user interface last
	if err = o.web.Serve(o.errc); err != nil {
		return err
	}

	log.Info().Msg("otterdb has started")
	if err = <-o.errc; err != nil {
		log.WithLevel(zerolog.FatalLevel).Err(err).Msg("otterdb has crashed")
		return err
	}
	return nil
}

func (o *OtterDB) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down otterdb")

	// Shutdown services in reverse order
	if serr := o.web.Shutdown(); serr != nil {
		err = errors.Join(err, serr)
	}

	if serr := o.server.Shutdown(); serr != nil {
		err = errors.Join(err, serr)
	}

	if serr := o.replica.Shutdown(); serr != nil {
		err = errors.Join(err, serr)
	}

	log.Debug().Msg("all otterdb services have shutdown")
	return err
}
