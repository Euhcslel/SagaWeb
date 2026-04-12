package sessions

import (
	"time"

	"gorm.io/gorm"
	"project/internal/domain/users"
)

const TableNameSession = "sessions"

type Session struct {
	UserID    int64     `gorm:"column:user_id;primaryKey"`
	Token     string    `gorm:"column:token;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`

	User users.User `gorm:"foreignKey:UserID;references:ID"`
}

func (*Session) TableName() string {
	return TableNameSession
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	tx.Where("user_id = ?", s.UserID).Delete(&Session{})
	s.ExpiresAt = time.Now().AddDate(0, 1, 0)
	return nil
}
