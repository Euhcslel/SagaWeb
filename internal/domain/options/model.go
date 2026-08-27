// Package options предоставляет модель для работы с дополнительными опциями.
// Данные о дополнительных опциях хранятся в таблице options.

package options

import (
	"github.com/shopspring/decimal"
)

const TableNameOption = "options"

type Option struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	ImagePath      string          `gorm:"column:image_path"`
}

func (*Option) TableName() string {
	return TableNameOption
}
