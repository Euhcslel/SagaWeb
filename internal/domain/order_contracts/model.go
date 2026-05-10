package order_contracts

import (
	"time"
)

const TableNameOrderContract = "order_contracts"

type OrderContract struct {
	OrderID        int64     `gorm:"column:order_id;primaryKey"`
	ContractNumber string    `gorm:"column:contract_number;primaryKey"`
	Path           string    `gorm:"column:path;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*OrderContract) TableName() string {
	return TableNameOrderContract
}
