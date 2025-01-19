package repository

import (
	"rest-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("session_token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) InvalidateSession(sessionID string) error {
	return r.db.Model(&models.Session{}).Where("id = ?", sessionID).Update("is_active", false).Error
}

func (r *sessionRepository) InvalidateAllUserSessions(userID string) error {
	return r.db.Model(&models.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
}
