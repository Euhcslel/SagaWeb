// Package order_appendices предоставляет модель для работы с приложениями к договору дилера.
// Данные о каждом приложении к договору хранятся в таблице order_appendices.

package order_appendices

import (
	"time"
)

const TableNameOrderAppendice = "order_appendices"

type OrderAppendix struct {
	OrderID         int64     `gorm:"column:order_id;primaryKey"`
	AppendiceNumber  string    `gorm:"column:appendice_number;primaryKey"`
	Path            string    `gorm:"column:path;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (*OrderAppendix) TableName() string {
	return TableNameOrderAppendice
}
