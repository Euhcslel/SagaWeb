package models

const TableNameLiftType = "lift_types"

type LiftType struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name            string `gorm:"column:name;not null" json:"name"`
	MinHeadroom int32 `gorm:"column:min_headroom;not null" json:"min_headroom"`
	MaxHeadroom int32 `gorm:"column:max_headroom;not null" json:"max_headroom"`
	WholesaleMarkup int32  `gorm:"column:wholesale_markup" json:"wholesale_markup"`
	RetailMarkup    int32  `gorm:"column:retail_markup" json:"retail_markup"`
}

func (*LiftType) TableName() string {
	return TableNameLiftType
}
