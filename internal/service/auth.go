package service

import (
	"errors"
	"time"

	apperrors "rest-api/internal/errors"
	"rest-api/internal/models"
	"rest-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   []byte
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) CreateSession(user *models.User, deviceInfo string) (*models.Session, error) {
	session := &models.Session{
		UserID:         user.ID,
		SessionToken:   uuid.New().String(),
		DeviceInfo:     deviceInfo,
		IsActive:       true,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		LastAccessedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) InvalidateSession(token string) error {
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return errors.New("session not found")
	}

	return s.sessionRepo.InvalidateSession(session.ID.String())
}

func (s *AuthService) InvalidateAllSessions(userID string) error {
	return s.sessionRepo.InvalidateAllUserSessions(userID)
}

func (s *AuthService) ValidateSession(token string) (*models.Session, error) {
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, apperrors.ErrSessionNotFound
	}

	if !session.IsActive {
		return nil, apperrors.ErrInvalidSession
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, apperrors.ErrSessionExpired
	}

	return session, nil
}

func (s *AuthService) Authenticate(email, password string) (*models.Session, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.CreateSession(user, "")
}
