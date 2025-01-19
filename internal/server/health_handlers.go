package server

import (
	"rest-api/internal/response"

	"github.com/gin-gonic/gin"
)

// @Summary Health check
// @Description Check if the server and database are healthy
// @Tags health
// @Produce json
// @Success 200 {object} response.HealthResponse
// @Router /health [get]
func (s *Server) handleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"version": s.cfg.Server.Version,
		})
	}
}
