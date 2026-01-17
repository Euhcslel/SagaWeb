package models

import (
	"time"

	"gorm.io/gorm"
)

const TableNameSession = "sessions"

type Session struct {
	UserID    int64     `gorm:"column:user_id;primaryKey" json:"user_id"`
	TokenHash string    `gorm:"column:token_hash;not null" json:"token_hash"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
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
