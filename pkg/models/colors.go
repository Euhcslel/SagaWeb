package models

const TableNameColor = "colors"

type Color struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Code string `gorm:"column:code;not null"`
	Hex  string `gorm:"column:hex;not null"`
}

func (*Color) TableName() string {
	return TableNameColor
}
