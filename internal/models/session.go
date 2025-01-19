package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	SessionToken   string    `gorm:"uniqueIndex;not null"`
	DeviceInfo     string
	IsActive       bool `gorm:"default:true"`
	ExpiresAt      time.Time
	LastAccessedAt time.Time
	CreatedAt      time.Time
	User           User `gorm:"foreignKey:UserID"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
