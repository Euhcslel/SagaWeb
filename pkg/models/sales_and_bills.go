package models

import (
	"time"
)

const TableNameSalesAndBill = "sales_and_bills"

type SalesAndBill struct {
	SaleID     int64     `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	BillNumber string    `gorm:"column:bill_number;primaryKey" json:"bill_number"`
	Path       string    `gorm:"column:path;not null" json:"path"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*SalesAndBill) TableName() string {
	return TableNameSalesAndBill
}
