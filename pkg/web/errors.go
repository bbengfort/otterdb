package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Renders the "not found page"
func (s *Server) NotFound(c *gin.Context) {
	c.Negotiate(http.StatusNotFound, gin.Negotiate{
		Offered: []string{binding.MIMEJSON, binding.MIMEHTML},
		Data: gin.H{
			"success": false,
			"error":   "not found",
		},
		HTMLName: "404.html",
	})
}

// Renders the "invalid action page"
func (s *Server) NotAllowed(c *gin.Context) {
	c.Negotiate(http.StatusMethodNotAllowed, gin.Negotiate{
		Offered: []string{binding.MIMEJSON, binding.MIMEHTML},
		Data: gin.H{
			"success": false,
			"error":   "method not allowed",
		},
		HTMLName: "405.html",
	})
}
