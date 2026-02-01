package models

import (
	"time"
)

const TableNameSalesAndBill = "sales_and_bills"

type SalesAndBill struct {
	SaleID     int64     `gorm:"column:sale_id;primaryKey"`
	BillNumber string    `gorm:"column:bill_number;primaryKey"`
	Path       string    `gorm:"column:path;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*SalesAndBill) TableName() string {
	return TableNameSalesAndBill
}
