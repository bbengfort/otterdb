/*
Package replica manages replication between peers in an otterdb cluster.
*/
package replica

import (
	"fmt"
	"net"
	"time"

	"github.com/bbengfort/otterdb/pkg/config"
	"github.com/bbengfort/otterdb/pkg/grpc/health/v1"
	"github.com/bbengfort/otterdb/pkg/replica/raft/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Replica struct {
	health.ProbeServer
	raft.UnimplementedRaftServer

	conf    config.ReplicaConfig
	srv     *grpc.Server
	started time.Time
}

func New(conf config.ReplicaConfig) (r *Replica, err error) {
	// Must supply a valid configuration.
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	r = &Replica{conf: conf}

	// Prepare to receive gRPC requests and configure RPCs
	opts := make([]grpc.ServerOption, 0, 4)
	// opts = append(opts, grpc.ChainUnaryInterceptor(s.UnaryInterceptors()...))
	// opts = append(opts, grpc.ChainStreamInterceptor(s.StreamInterceptors()...))
	r.srv = grpc.NewServer(opts...)

	// Initialize the gRPC services
	raft.RegisterRaftServer(r.srv, r)
	health.RegisterHealthServer(r.srv, r)

	// Set the server to a not serving state
	r.NotHealthy()

	return r, nil
}

func (r *Replica) Serve(errc chan<- error) (err error) {
	if !r.conf.Enabled {
		log.Warn().Bool("enabled", r.conf.Enabled).Msg("otterdb replication is disabled")
		return nil
	}

	// Listen for TCP requests (other sockets such as bufconn for tests should use Run)
	var sock net.Listener
	if sock, err = net.Listen("tcp", r.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %w", r.conf.BindAddr, err)
	}

	// Run the server on the opened socket
	go r.Run(errc, sock)

	// Now that the server is running mark healthy and set the start time to track uptime
	r.started = time.Now()
	r.Healthy()

	log.Info().Str("listen", r.conf.BindAddr).Msg("otterdb replica server started")
	return nil
}

// Run the gRPC server on the specified socket. This method can be used to serve TCP
// requests or to connect to a bufconn for testing purposes. This method blocks while
// the server is running so it should be run in a go routine.
func (r *Replica) Run(errc chan<- error, sock net.Listener) {
	defer sock.Close()
	if err := r.srv.Serve(sock); err != nil {
		errc <- err
	}
}

func (r *Replica) Shutdown() (err error) {
	// If the server is not enabled, skip shutdown
	if !r.conf.Enabled {
		return nil
	}

	// Set the server to a not serving state
	log.Debug().Msg("gracefully shutting down otterdb replica server")
	r.NotHealthy()

	r.srv.GracefulStop()
	return nil
}
