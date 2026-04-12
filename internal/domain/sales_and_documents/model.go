package sales_and_documents

import (
	"time"
)

const TableNameSalesAndDocument = "sales_and_documents"

type SalesAndDocument struct {
	SaleID    int64     `gorm:"column:sale_id;primaryKey"`
	Name      string    `gorm:"column:name;primaryKey"`
	Path      string    `gorm:"column:path;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*SalesAndDocument) TableName() string {
	return TableNameSalesAndDocument
}
