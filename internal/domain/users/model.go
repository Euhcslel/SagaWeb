package users

import (
	"project/internal/domain/roles"
	"time"
)

const TableNameUser = "users"

type User struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true"`
	Username     string    `gorm:"column:username;not null;uniqueIndex"`
	Fullname     string    `gorm:"column:fullname;not null;"`
	Email        string    `gorm:"column:email"`
	PhoneNumber  string    `gorm:"column:phone_number;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	PasswordHash []byte    `gorm:"column:password_hash;not null"`
	RoleID       int64     `gorm:"column:role_id;not null"`

	Role roles.Role `gorm:"foreignKey:RoleID;references:ID"`
}

func (*User) TableName() string {
	return TableNameUser
}
