/*
Package server implements the query server for making database queries to otterdb.
*/
package server

import "github.com/bbengfort/otterdb/pkg/config"

type Server struct {
	conf config.ServerConfig
}

func New(conf config.ServerConfig) (srv *Server, err error) {
	return &Server{conf: conf}, nil
}

func (s *Server) Serve(errc chan<- error) (err error) {
	return nil
}

func (s *Server) Shutdown() (err error) {
	return nil
}
