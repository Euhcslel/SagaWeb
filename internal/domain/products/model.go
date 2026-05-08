package products

import (
	"github.com/shopspring/decimal"
)

const TableNameProduct = "products"

type Product struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
}

func (*Product) TableName() string {
	return TableNameProduct
}
