package units

const TableNameUnit = "units"

type Unit struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name string `gorm:"column:name;not null"`
}

func (*Unit) TableName() string {
	return TableNameUnit
}
