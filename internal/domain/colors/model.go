// Package colors предоставляет модель для работы с цветами.
// Цвета хранятся в таблице colors.

package colors

const TableNameColor = "colors"

type Color struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Code string `gorm:"column:code;not null"`
	Hex  string `gorm:"column:hex;not null"`
}

func (*Color) TableName() string {
	return TableNameColor
}
