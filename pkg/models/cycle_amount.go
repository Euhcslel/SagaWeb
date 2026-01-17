package models

const TableNameCycleAmount = "cycle_amount"

type CycleAmount struct {
	ID              int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Amount          int32 `gorm:"column:amount;not null" json:"amount"`
	WholesaleMarkup int32 `gorm:"column:wholesale_markup" json:"wholesale_markup"`
	RetailMarkup    int32 `gorm:"column:retail_markup" json:"retail_markup"`
}

func (*CycleAmount) TableName() string {
	return TableNameCycleAmount
}
