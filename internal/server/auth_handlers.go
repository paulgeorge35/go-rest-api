package server

import (
	"fmt"
	"rest-api/internal/errors"
	"rest-api/internal/response"
	"rest-api/internal/service"
	"rest-api/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.RegisterRequest true "Register Request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /v1/register [post]
func (s *Server) handleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validator.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Find user by email
		_, err := s.userSvc.FindByEmail(req.Email)
		if err == nil {
			response.BadRequest(c, errors.ErrEmailAlreadyExists)
			return
		}

		// Create new user
		newUser, err := s.userSvc.Register(req.Email, req.Password, req.Name)
		if err != nil {
			s.logger.Error("failed to create user", err)
			response.InternalError(c, errors.ErrFailedToCreateUser)
			return
		}

		// Create session after registration
		session, err := s.authSvc.CreateSession(newUser, c.GetHeader("User-Agent"))
		if err != nil {
			s.logger.Error("failed to create session", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		response.SuccessWithMessage(c, "registration successful", gin.H{
			"token": session.SessionToken,
		})
	}
}

// @Summary Login user
// @Description Authenticate user and return session token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.LoginRequest true "Login Request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /v1/login [post]
func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validator.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Validate request
		if err := s.validator.ValidateLoginRequest(&req); err != nil {
			response.BadRequest(c, err)
			return
		}

		// Authenticate user
		session, err := s.authSvc.Authenticate(req.Email, req.Password)
		if err != nil {
			s.logger.Error("failed to authenticate user", err)
			response.Unauthorized(c, errors.ErrInvalidEmailOrPass)
			return
		}

		response.SuccessWithMessage(c, "login successful", gin.H{
			"token": session.SessionToken,
		})
	}
}

// @Summary Request password reset
// @Description Send password reset link to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.ForgotPasswordRequest true "Email Request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/forgot-password [post]
func (s *Server) handleForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validator.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Find user by email
		user, err := s.userSvc.FindByEmail(req.Email)
		if err != nil {
			// Don't reveal if email exists
			response.SuccessWithMessage(c, "if the email exists, a password reset link will be sent", nil)
			return
		}

		// Generate reset token
		token, err := s.tokenSvc.GenerateToken(user.ID.String(), service.TokenTypeReset)
		if err != nil {
			s.logger.Error("failed to generate reset token", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		// Create reset link and send email
		resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.Server.BaseURL, token)
		emailBody := fmt.Sprintf(`
			<h1>Password Reset</h1>
			<p>Click the link below to reset your password:</p>
			<a href="%s">Reset Password</a>
			<p>This link will expire in 15 minutes.</p>
		`, resetLink)

		if err := s.emailSvc.SendEmail(user.Email, "Password Reset", emailBody); err != nil {
			s.logger.Error("failed to send reset email", err)
			response.InternalError(c, errors.ErrFailedToSendEmail)
			return
		}

		response.SuccessWithMessage(c, "if the email exists, a password reset link will be sent", nil)
	}
}

// @Summary Request magic link login
// @Description Send magic link to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.MagicLinkRequest true "Email Request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/magic-link-login [post]
func (s *Server) handleMagicLinkLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validator.MagicLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Find user by email
		user, err := s.userSvc.FindByEmail(req.Email)
		if err != nil {
			// Don't reveal if email exists
			response.SuccessWithMessage(c, "if the email exists, a magic link will be sent", nil)
			return
		}

		// Generate magic link token
		token, err := s.tokenSvc.GenerateToken(user.ID.String(), service.TokenTypeMagicLink)
		if err != nil {
			s.logger.Error("failed to generate magic link token", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		// Create magic link and send email
		magicLink := fmt.Sprintf("%s/magic-login?token=%s", s.cfg.Server.BaseURL, token)
		emailBody := fmt.Sprintf(`
			<h1>Magic Link Login</h1>
			<p>Click the link below to log in:</p>
			<a href="%s">Log In</a>
			<p>This link will expire in 15 minutes.</p>
		`, magicLink)

		if err := s.emailSvc.SendEmail(user.Email, "Magic Link Login", emailBody); err != nil {
			s.logger.Error("failed to send magic link email", err)
			response.InternalError(c, errors.ErrFailedToSendEmail)
			return
		}

		response.SuccessWithMessage(c, "if the email exists, a magic link will be sent", nil)
	}
}

// @Summary Get Google OAuth URL
// @Description Get Google OAuth URL for authentication
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/oauth/google [get]
func (s *Server) handleGoogleOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate state token to prevent CSRF
		state := uuid.New().String()

		// Store state in cookie
		c.SetCookie("oauth_state", state, 600, "/", "", false, true)

		// Get authorization URL
		authURL := s.oauthSvc.GetAuthURL(state)

		response.Success(c, gin.H{
			"redirect_url": authURL,
		})
	}
}

// @Summary Google OAuth Callback
// @Description Handle Google OAuth callback and create session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/oauth/google/callback [get]
func (s *Server) handleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify state token
		state, err := c.Cookie("oauth_state")
		if err != nil || state != c.Query("state") {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Clear state cookie
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)

		// Exchange code for token and get user info
		code := c.Query("code")
		session, err := s.oauthSvc.HandleCallback(code)
		if err != nil {
			s.logger.Error("failed to handle oauth callback", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		response.SuccessWithMessage(c, "login successful", gin.H{
			"token": session.SessionToken,
		})
	}
}

// @Summary Reset password
// @Description Reset user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/reset-password [post]
func (s *Server) handleResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validator.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Validate token
		tokenRecord, err := s.tokenSvc.ValidateToken(req.Token, service.TokenTypeReset)
		if err != nil {
			response.BadRequest(c, errors.ErrInvalidToken)
			return
		}

		// Update user's password
		if err := s.userSvc.UpdatePassword(tokenRecord.UserID, req.Password); err != nil {
			s.logger.Error("failed to update password", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		// Invalidate the token
		if err := s.tokenSvc.InvalidateToken(req.Token); err != nil {
			s.logger.Error("failed to invalidate token", err)
		}

		response.SuccessWithMessage(c, "password updated successfully", nil)
	}
}

func (s *Server) handleMagicLinkVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			response.BadRequest(c, errors.ErrInvalidRequest)
			return
		}

		// Validate token
		tokenRecord, err := s.tokenSvc.ValidateToken(token, service.TokenTypeMagicLink)
		if err != nil {
			response.BadRequest(c, errors.ErrInvalidToken)
			return
		}

		// Get user
		user, err := s.userSvc.GetByID(tokenRecord.UserID)
		if err != nil {
			s.logger.Error("failed to get user", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		// Create session
		deviceInfo := c.GetHeader("User-Agent")
		session, err := s.authSvc.CreateSession(user, deviceInfo)
		if err != nil {
			s.logger.Error("failed to create session", err)
			response.InternalError(c, errors.ErrInvalidRequest)
			return
		}

		// Invalidate the token
		if err := s.tokenSvc.InvalidateToken(token); err != nil {
			s.logger.Error("failed to invalidate token", err)
		}

		response.SuccessWithMessage(c, "login successful", gin.H{
			"token": session.SessionToken,
		})
	}
}
