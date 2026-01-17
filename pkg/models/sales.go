package models

import (
	"time"
)

const TableNameSale = "sales"

type Sale struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ClientID    int64     `gorm:"column:client_id" json:"client"`
	ManagerID   int64     `gorm:"column:manager_id" json:"manager"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

    Client  User `gorm:"foreignKey:ClientID;references:ID"`
    Manager User `gorm:"foreignKey:ManagerID;references:ID"`
}

func (*Sale) TableName() string {
	return TableNameSale
}
