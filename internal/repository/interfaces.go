package repository

import "rest-api/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	UpdatePassword(userID string, hashedPassword string) error
}

type SessionRepository interface {
	Create(session *models.Session) error
	FindByToken(token string) (*models.Session, error)
	InvalidateSession(sessionID string) error
	InvalidateAllUserSessions(userID string) error
}

type TokenRepository interface {
	Create(token *models.Token) error
	FindByToken(token string) (*models.Token, error)
	InvalidateToken(token string) error
}
