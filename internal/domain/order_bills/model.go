package order_bills

import (
	"time"
)

const TableNameOrderBill = "order_bills"

type OrderBill struct {
	OrderID    int64     `gorm:"column:order_id;primaryKey"`
	BillNumber string    `gorm:"column:bill_number;primaryKey"`
	Path       string    `gorm:"column:path;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*OrderBill) TableName() string {
	return TableNameOrderBill
}
