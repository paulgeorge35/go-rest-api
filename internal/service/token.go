package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"rest-api/internal/models"
	"rest-api/internal/repository"
)

type TokenType string

const (
	TokenTypeReset     TokenType = "reset"
	TokenTypeMagicLink TokenType = "magic_link"
)

type TokenService struct {
	tokenRepo repository.TokenRepository
}

func NewTokenService(tokenRepo repository.TokenRepository) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
	}
}

func (s *TokenService) GenerateToken(userID string, tokenType TokenType) (string, error) {
	// Generate random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// Create token record
	tokenRecord := &models.Token{
		UserID:    userID,
		Token:     token,
		Type:      string(tokenType),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.tokenRepo.Create(tokenRecord); err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenService) ValidateToken(token string, tokenType TokenType) (*models.Token, error) {
	tokenRecord, err := s.tokenRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if tokenRecord.Type != string(tokenType) {
		return nil, errors.New("invalid token type")
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return tokenRecord, nil
}

func (s *TokenService) InvalidateToken(token string) error {
	return s.tokenRepo.InvalidateToken(token)
}
