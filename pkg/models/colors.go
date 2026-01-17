package models

const TableNameColor = "colors"

type Color struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Code string `gorm:"column:code;not null" json:"code"`
	Hex  string `gorm:"column:hex;not null" json:"hex"`
}

func (*Color) TableName() string {
	return TableNameColor
}
