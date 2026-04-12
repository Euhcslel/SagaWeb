package products

import (
	"github.com/shopspring/decimal"
	"project/internal/domain/units"
)

const TableNameProduct = "products"

type Product struct {
	ID             int32           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	UnitID         int32           `gorm:"column:unit_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`

	Unit units.Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*Product) TableName() string {
	return TableNameProduct
}
