package models

const TableNameLiftType = "lift_types"

type LiftType struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name            string `gorm:"column:name;not null"`
	MinHeadroom     int32  `gorm:"column:min_headroom;not null"`
	MaxHeadroom     int32  `gorm:"column:max_headroom;not null"`
	WholesaleMarkup int32  `gorm:"column:wholesale_markup"`
	RetailMarkup    int32  `gorm:"column:retail_markup"`
}

func (*LiftType) TableName() string {
	return TableNameLiftType
}
