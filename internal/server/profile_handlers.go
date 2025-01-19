package server

import (
	"rest-api/internal/errors"
	"rest-api/internal/models"
	"rest-api/internal/response"

	"github.com/gin-gonic/gin"
)

// @Summary Get user profile
// @Description Get the current user's profile
// @Tags profile
// @Security Bearer
// @Produce json
// @Success 200 {object} response.UserResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /v1/profile [get]
func (s *Server) handleGetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, errors.ErrUnauthorized)
			return
		}

		response.Success(c, user.(*models.User))
	}
}

// @Summary Logout user
// @Description Invalidate current session
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /v1/logout [get]
func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, errors.ErrUnauthorized)
			return
		}

		if err := s.authSvc.InvalidateSession(token); err != nil {
			s.logger.Error("failed to invalidate session", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		response.SuccessWithMessage(c, "logged out successfully", nil)
	}
}

// @Summary Invalidate all sessions
// @Description Invalidate all user sessions
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /v1/invalidate-sessions [post]
func (s *Server) handleInvalidateSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, errors.ErrUnauthorized)
			return
		}

		if err := s.authSvc.InvalidateAllSessions(user.(*models.User).ID.String()); err != nil {
			s.logger.Error("failed to invalidate sessions", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		response.SuccessWithMessage(c, "all sessions invalidated", nil)
	}
}
