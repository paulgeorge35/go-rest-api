package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string    `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	Type      string    `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
