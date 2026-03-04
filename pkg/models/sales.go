package models

import (
	"time"
)

const TableNameSale = "sales"

type Sale struct {
	ID        int64       `gorm:"column:id;primaryKey;autoIncrement:true"`
	ClientID  *int64      `gorm:"column:client_id;null;default:NULL"`
	ManagerID int64       `gorm:"column:manager_id;not null"`
	CreatedAt time.Time   `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	Status    OrderStatus `gorm:"type:order_status;not null;default:'new'"`

	Client  *User `gorm:"foreignKey:ClientID;references:ID"`
	Manager User  `gorm:"foreignKey:ManagerID;references:ID"`
}

func (*Sale) TableName() string {
	return TableNameSale
}
