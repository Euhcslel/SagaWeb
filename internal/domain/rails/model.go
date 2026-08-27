// Package rails предоставляет модель для работы с направляющими.
// Данные о направляющих хранятся в таблице rails.

package rails

import (
	"github.com/shopspring/decimal"
)

const TableNameRail = "rails"

type Rail struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	Specifications *string         `gorm:"column:specifications"`
	ImagePath      string          `gorm:"column:image_path"`
}

func (*Rail) TableName() string {
	return TableNameRail
}
