package models

import (
	"time"

	"gorm.io/gorm"
)

const TableNameSession = "sessions"

type Session struct {
	UserID    int64     `gorm:"column:user_id;primaryKey"`
	TokenHash [32]byte    `gorm:"column:token_hash;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (*Session) TableName() string {
	return TableNameSession
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	tx.Where("user_id = ?", s.UserID).Delete(&Session{})
	s.ExpiresAt = time.Now().AddDate(0, 1, 0)
	return nil
}
