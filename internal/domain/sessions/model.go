// Package sessions предоставляет модель для работы с сессиями.
// Данные о сессиях хранятся в таблице sessions.

package sessions

import (
	"time"

	"gorm.io/gorm"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
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

// Функция-триггер, срабатывающая перед созданием новой сессии
// Удаляет старую сессию и устанавливает дату истечения новой сессии
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if err := tx.Where("user_id = ?", s.UserID).Delete(&Session{}).Error; err != nil {
		return err
	}
	s.ExpiresAt = time.Now().AddDate(0, 1, 0)
	return nil
}
