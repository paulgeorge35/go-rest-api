package repository

import (
	"rest-api/internal/models"
	"time"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) FindByToken(token string) (*models.Token, error) {
	var t models.Token
	err := r.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tokenRepository) InvalidateToken(token string) error {
	return r.db.Model(&models.Token{}).Where("token = ?", token).Update("used", true).Error
}
