package web

import (
	"net/http"
	"time"

	"github.com/bbengfort/otterdb/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	serverStatusOK          = "ok"
	serverStatusNotReady    = "not ready"
	serverStatusUnhealthy   = "unhealthy"
	serverStatusMaintenance = "maintenance"
)

// Status reports the version and uptime of the server
func (s *Server) Status(c *gin.Context) {
	var state string
	s.RLock()
	switch {
	case s.healthy && s.ready:
		state = serverStatusOK
	case s.healthy && !s.ready:
		state = serverStatusNotReady
	case !s.healthy:
		state = serverStatusUnhealthy
	}
	s.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"status":  state,
		"version": pkg.Version(),
		"uptime":  time.Since(s.started).String(),
	})
}

// If the server is in maintenance mode, aborts the current request and renders the
// maintenance mode page instead. Returns nil if not in maintenance mode.
func (s *Server) Maintenance() gin.HandlerFunc {
	if s.conf.Maintenance {
		return func(c *gin.Context) {
			c.Negotiate(http.StatusServiceUnavailable, gin.Negotiate{
				Offered: []string{binding.MIMEJSON, binding.MIMEHTML},
				Data: gin.H{
					"status":  serverStatusMaintenance,
					"version": pkg.Version(),
					"uptime":  time.Since(s.started).String(),
				},
				HTMLName: "maintenance.html",
			})
			c.Abort()
		}
	}
	return nil
}

// Healthz is used to alert k8s to the health/liveness status of the server.
func (s *Server) Healthz(c *gin.Context) {
	s.RLock()
	healthy := s.healthy
	s.RUnlock()

	if !healthy {
		c.Data(http.StatusServiceUnavailable, "text/plain", []byte(serverStatusUnhealthy))
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(serverStatusOK))
}

// Readyz is used to alert k8s to the readiness status of the server.
func (s *Server) Readyz(c *gin.Context) {
	s.RLock()
	ready := s.ready
	s.RUnlock()

	if !ready {
		c.Data(http.StatusServiceUnavailable, "text/plain", []byte(serverStatusNotReady))
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(serverStatusOK))
}
