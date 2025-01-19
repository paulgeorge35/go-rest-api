package server

import (
	"net/http"
	"rest-api/internal/errors"
	"rest-api/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Join multiple origins with comma if needed
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(s.cfg.CORS.AllowedOrigins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func (s *Server) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := s.rateLimiter.GetLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")

		// Validate session
		session, err := s.authSvc.ValidateSession(token)
		if err != nil {
			response.Unauthorized(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Get user
		user, err := s.userSvc.GetByID(session.UserID.String())
		if err != nil {
			response.Unauthorized(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}
