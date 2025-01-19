package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         string
		BaseURL      string
		Env          string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		Version      string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	JWT struct {
		Secret string
	}
	OAuth struct {
		Google struct {
			ClientID     string
			ClientSecret string
			RedirectURL  string
		}
	}
	CORS struct {
		AllowedOrigins []string
	}
	Email struct {
		Host     string
		Port     string
		Username string
		Password string
		From     string
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}

	// Server
	cfg.Server.Port = os.Getenv("PORT")
	cfg.Server.BaseURL = os.Getenv("BASE_URL")
	cfg.Server.Env = os.Getenv("ENV")
	cfg.Server.Version = os.Getenv("VERSION")

	// Database
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Name = os.Getenv("DB_NAME")

	// JWT
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	// OAuth
	cfg.OAuth.Google.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.OAuth.Google.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.OAuth.Google.RedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")

	// CORS
	cfg.CORS.AllowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	// Email
	cfg.Email.Host = os.Getenv("EMAIL_HOST")
	cfg.Email.Port = os.Getenv("EMAIL_PORT")
	cfg.Email.Username = os.Getenv("EMAIL_USERNAME")
	cfg.Email.Password = os.Getenv("EMAIL_PASSWORD")
	cfg.Email.From = os.Getenv("EMAIL_FROM")

	// Set default timeouts
	cfg.Server.ReadTimeout = 15 * time.Second
	cfg.Server.WriteTimeout = 15 * time.Second

	return cfg, nil
}
