/*
Package web provides an implementation of a web UI that can be enabled to access and
manage a single replica or the cluster of replicas if that replica is the leader.
*/
package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bbengfort/otterdb/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	sync.RWMutex
	conf    config.WebConfig
	srv     *http.Server
	router  *gin.Engine
	url     *url.URL
	started time.Time
	healthy bool
	ready   bool
}

func New(conf config.WebConfig) (srv *Server, err error) {
	// Must supply a valid configuration.
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	srv = &Server{conf: conf}

	// If not enabled, return just the server stub
	if !conf.Enabled {
		return srv, nil
	}

	// Configure the gin router when enabled
	srv.router = gin.New()
	srv.router.RedirectTrailingSlash = true
	srv.router.RedirectFixedPath = false
	srv.router.HandleMethodNotAllowed = true
	srv.router.ForwardedByClientIP = true
	srv.router.UseRawPath = false
	srv.router.UnescapePathValues = true
	if err = srv.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server if enabled
	srv.srv = &http.Server{
		Addr:              srv.conf.BindAddr,
		Handler:           srv.router,
		ErrorLog:          nil,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return srv, nil
}

func (s *Server) Serve(errc chan<- error) (err error) {
	if !s.conf.Enabled {
		log.Warn().Bool("enabled", s.conf.Enabled).Msg("otterdb web ui is disabled")
		return nil
	}

	// Create a socket to listen on and infer the final URL.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.srv.Addr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %w", s.srv.Addr, err)
	}

	s.setURL(sock.Addr())
	s.SetStatus(true, true)
	s.started = time.Now()

	// Listen for HTTP requests and handle them
	go func() {
		// Make sure we don't use the external err to avoid data races.
		if serr := s.serve(sock); !errors.Is(serr, http.ErrServerClosed) {
			errc <- serr
		}
	}()

	log.Info().Str("url", s.URL()).Msg("otterdb user interface started")
	return nil
}

// ServeTLS if a tls configuration is provided, otherwise Serve.
func (s *Server) serve(sock net.Listener) error {
	if s.srv.TLSConfig != nil {
		return s.srv.ServeTLS(sock, "", "")
	}
	return s.srv.Serve(sock)
}

func (s *Server) Shutdown() (err error) {
	// If the server is not enabled, skip shutdown
	if !s.conf.Enabled {
		return nil
	}

	log.Debug().Msg("gracefully shutting down web user interface server")
	s.SetStatus(false, false)

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// SetStatus sets the health and ready status on the server, modifying the behavior of
// the kubernetes probe responses.
func (s *Server) SetStatus(health, ready bool) {
	s.Lock()
	s.healthy = health
	s.ready = ready
	s.Unlock()
	log.Debug().Bool("health", health).Bool("ready", ready).Msg("server status set")
}

// URL returns the endpoint of the server as determined by the configuration and the
// socket address and port (if specified).
func (s *Server) URL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url.String()
}

func (s *Server) setURL(addr net.Addr) {
	s.Lock()
	defer s.Unlock()

	s.url = &url.URL{
		Scheme: "http",
		Host:   addr.String(),
	}

	if s.srv.TLSConfig != nil {
		s.url.Scheme = "https"
	}

	if tcp, ok := addr.(*net.TCPAddr); ok && tcp.IP.IsUnspecified() {
		s.url.Host = fmt.Sprintf("127.0.0.1:%d", tcp.Port)
	}
}

// Debug returns a server that uses the specified http server instead of creating one.
// This function is primarily used to create test servers easily.
func Debug(conf config.WebConfig, srv *http.Server) (s *Server, err error) {
	if s, err = New(conf); err != nil {
		return nil, err
	}

	// Replace the http server with the one specified
	s.srv = nil
	s.srv = srv
	s.srv.Handler = s.router
	return s, nil
}
