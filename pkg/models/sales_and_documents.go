package models

import (
	"time"
)

const TableNameSalesAndDocument = "sales_and_documents"

type SalesAndDocument struct {
	SaleID    int64     `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	Name      string    `gorm:"column:name;primaryKey" json:"name"`
	Path      string    `gorm:"column:path;not null" json:"path"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*SalesAndDocument) TableName() string {
	return TableNameSalesAndDocument
}
