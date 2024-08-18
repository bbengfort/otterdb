/*
Package web provides an implementation of a web UI that can be enabled to access and
manage a single replica or the cluster of replicas if that replica is the leader.
*/
package web

import "github.com/bbengfort/otterdb/pkg/config"

type Server struct {
	conf config.WebConfig
}

func New(conf config.WebConfig) (srv *Server, err error) {
	return &Server{conf: conf}, nil
}

func (s *Server) Serve(errc chan<- error) (err error) {
	return nil
}

func (s *Server) Shutdown() (err error) {
	return nil
}
