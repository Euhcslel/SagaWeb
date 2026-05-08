package users

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"time"
)

const TableNameUser = "users"

type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true"`
	Fullname     string     `gorm:"column:fullname;not null;"`
	Email        string     `gorm:"column:email"`
	PhoneNumber  string     `gorm:"column:phone_number;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	PasswordHash []byte     `gorm:"column:password_hash;not null"`
	Role         enums.Role `gorm:"column:role;not null"`
}

func (*User) TableName() string {
	return TableNameUser
}
