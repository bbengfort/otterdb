/*
Package server implements the query server for making database queries to otterdb.
*/
package server

import (
	"fmt"
	"net"
	"time"

	"github.com/bbengfort/otterdb/pkg/config"
	"github.com/bbengfort/otterdb/pkg/grpc/health/v1"
	"github.com/bbengfort/otterdb/pkg/server/api/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	health.ProbeServer
	api.UnimplementedOtterServer

	conf    config.ServerConfig
	srv     *grpc.Server
	started time.Time
}

func New(conf config.ServerConfig) (s *Server, err error) {
	// Must supply a valid configuration.
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	s = &Server{conf: conf}

	// Prepare to receive gRPC requests and configure RPCs
	opts := make([]grpc.ServerOption, 0, 4)
	// opts = append(opts, grpc.ChainUnaryInterceptor(s.UnaryInterceptors()...))
	// opts = append(opts, grpc.ChainStreamInterceptor(s.StreamInterceptors()...))
	s.srv = grpc.NewServer(opts...)

	// Initialize the gRPC services
	api.RegisterOtterServer(s.srv, s)
	health.RegisterHealthServer(s.srv, s)

	// Set the server to a not serving state
	s.NotHealthy()

	return s, nil
}

func (s *Server) Serve(errc chan<- error) (err error) {
	if !s.conf.Enabled {
		log.Warn().Bool("enabled", s.conf.Enabled).Msg("otterdb database server is disabled")
		return nil
	}

	// Listen for TCP requests (other sockets such as bufconn for tests should use Run)
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %w", s.conf.BindAddr, err)
	}

	// Run the server on the opened socket
	go s.Run(errc, sock)

	// Now that the server is running mark healthy and set the start time to track uptime
	s.started = time.Now()
	s.Healthy()

	log.Info().Str("listen", s.conf.BindAddr).Msg("otterdb database server started")
	return nil
}

// Run the gRPC server on the specified socket. This method can be used to serve TCP
// requests or to connect to a bufconn for testing purposes. This method blocks while
// the server is running so it should be run in a go routine.
func (s *Server) Run(errc chan<- error, sock net.Listener) {
	defer sock.Close()
	if err := s.srv.Serve(sock); err != nil {
		errc <- err
	}
}

func (s *Server) Shutdown() (err error) {
	// If the server is not enabled, skip shutdown
	if !s.conf.Enabled {
		return nil
	}

	// Set the server to a not serving state
	log.Debug().Msg("gracefully shutting down otterdb database server")
	s.NotHealthy()

	s.srv.GracefulStop()
	return nil
}
