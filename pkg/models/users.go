package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const TableNameUser = "users"

type User struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username    string    `gorm:"column:username;not null;uniqueIndex" json:"username"`
	Fullname    string    `gorm:"column:fullname;not null;" json:"fullname"`
	Email       string    `gorm:"column:email" json:"email"`
	PhoneNumber string    `gorm:"column:phone_number;not null" json:"phone_number"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	Password    string    `gorm:"column:password;not null;uniqueIndex" json:"password"`
	RoleID      int64     `gorm:"column:role_id;not null" json:"role_id"`

	Role Role `gorm:"foreignKey:RoleID;references:ID"`
}

func (*User) TableName() string {
	return TableNameUser
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Ошибка при попытке хэширования пароля")
		return err
	}

	u.Password = string(passwordHash)
	return nil
}
