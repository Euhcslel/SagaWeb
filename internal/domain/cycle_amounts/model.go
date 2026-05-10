package cycle_amounts

import "github.com/shopspring/decimal"

const TableNameCycleAmount = "cycle_amounts"

type CycleAmount struct {
	ID              int64           `gorm:"column:id;primaryKey;autoIncrement:true"`
	Amount          string          `gorm:"column:amount;not null"`
	WholesaleMarkup decimal.Decimal `gorm:"column:wholesale_markup"`
	RetailMarkup    decimal.Decimal `gorm:"column:retail_markup"`
}

func (*CycleAmount) TableName() string {
	return TableNameCycleAmount
}
