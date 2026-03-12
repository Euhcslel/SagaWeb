package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameCycleAmountMarkupHistory = "cycle_amount_markup_history"

type CycleAmountMarkupHistory struct {
	ID              int64           `gorm:"column:id;primaryKey"`
	CycleAmountID   int64           `gorm:"column:cycle_amount_id"`
	WholesaleMarkup decimal.Decimal `gorm:"column:wholesale_markup"`
	RetailMarkup    decimal.Decimal `gorm:"column:retail_markup"`
	SetAt           time.Time       `gorm:"column:set_at"`
}

func (*CycleAmountMarkupHistory) TableName() string {
	return TableNameCycleAmountMarkupHistory
}
