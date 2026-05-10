package order_documents

import (
	"time"
)

const TableNameOrderDocument = "order_documents"

type OrderDocument struct {
	OrderID   int64     `gorm:"column:order_id;primaryKey"`
	Name      string    `gorm:"column:name;primaryKey"`
	Path      string    `gorm:"column:path;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*OrderDocument) TableName() string {
	return TableNameOrderDocument
}
