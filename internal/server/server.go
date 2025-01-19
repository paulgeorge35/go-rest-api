package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rest-api/internal/config"
	"rest-api/internal/middleware"
	"rest-api/internal/repository"
	"rest-api/internal/service"
	"rest-api/pkg/email"
	"rest-api/pkg/logger"
	"rest-api/pkg/validator"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type Server struct {
	cfg         *config.Config
	logger      *logger.Logger
	router      *gin.Engine
	httpSrv     *http.Server
	validator   *validator.Validator
	authSvc     *service.AuthService
	userSvc     *service.UserService
	tokenSvc    *service.TokenService
	emailSvc    *email.EmailService
	rateLimiter *middleware.IPRateLimiter
	oauthSvc    *service.OAuthService
	db          *gorm.DB
}

func New(cfg *config.Config, logger *logger.Logger, db *gorm.DB) *Server {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize services
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, sessionRepo, cfg.JWT.Secret)
	tokenSvc := service.NewTokenService(tokenRepo)

	// Convert port string to int
	port, _ := strconv.Atoi(cfg.Email.Port)

	emailSvc := email.NewEmailService(email.Config{
		Host:     cfg.Email.Host,
		Port:     port, // Now using the converted int value
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
		From:     cfg.Email.From,
	})

	// Initialize OAuth service
	oauthSvc := service.NewOAuthService(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Google.RedirectURL,
		userSvc,
		authSvc,
	)

	return &Server{
		cfg:         cfg,
		logger:      logger,
		router:      gin.Default(),
		validator:   validator.NewValidator(),
		userSvc:     userSvc,
		authSvc:     authSvc,
		tokenSvc:    tokenSvc,
		emailSvc:    emailSvc,
		rateLimiter: middleware.NewIPRateLimiter(rate.Limit(1), 5),
		oauthSvc:    oauthSvc,
		db:          db,
	}
}

func (s *Server) Start() error {
	// Setup routes
	s.setupRoutes()

	// Configure HTTP server
	s.httpSrv = &http.Server{
		Addr:         ":" + s.cfg.Server.Port,
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}

	// Start server
	s.logger.Info(fmt.Sprintf("Starting server on port %s", s.cfg.Server.Port))
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

func (s *Server) setupRoutes() {
	// Add Swagger
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (no versioning)
	s.router.GET("/api/health", s.handleHealthCheck())

	// Add middleware
	s.router.Use(gin.Recovery())
	s.router.Use(s.corsMiddleware())
	s.router.Use(s.rateLimitMiddleware())

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", s.handleRegister())
		v1.POST("/login", s.handleLogin())
		v1.POST("/forgot-password", s.handleForgotPassword())
		v1.POST("/magic-link-login", s.handleMagicLinkLogin())
		v1.GET("/oauth/google", s.handleGoogleOAuth())
		v1.GET("/oauth/google/callback", s.handleGoogleCallback())
		v1.POST("/reset-password", s.handleResetPassword())
		v1.GET("/verify-magic-link", s.handleMagicLinkVerify())

		// Protected routes
		protected := v1.Group("/")
		protected.Use(s.authMiddleware())
		{
			protected.GET("/profile", s.handleGetProfile())
			protected.GET("/logout", s.handleLogout())
			protected.POST("/invalidate-sessions", s.handleInvalidateSessions())
		}
	}
}
