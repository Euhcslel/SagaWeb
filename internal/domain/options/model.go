package options

import (
	"github.com/shopspring/decimal"
	"project/internal/domain/units"
)

const TableNameOption = "options"

type Option struct {
	ID             int32           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name           string          `gorm:"column:name;not null"`
	UnitID         int32           `gorm:"column:unit_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	ForSale        bool            `gorm:"column:for_sale;not null"`
	Condition      string          `gorm:"column:condition"`

	Unit units.Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*Option) TableName() string {
	return TableNameOption
}
