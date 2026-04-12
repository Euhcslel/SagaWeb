package lift_type_markup_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameLiftTypeMarkupHistory = "lift_type_markup_history"

type LiftTypeMarkupHistory struct {
	ID              int64           `gorm:"column:id;primaryKey"`
	LiftTypeID      int64           `gorm:"column:lift_type_id"`
	WholesaleMarkup decimal.Decimal `gorm:"column:wholesale_markup"`
	RetailMarkup    decimal.Decimal `gorm:"column:retail_markup"`
	SetAt           time.Time       `gorm:"column:set_at"`
}

func (*LiftTypeMarkupHistory) TableName() string {
	return TableNameLiftTypeMarkupHistory
}
