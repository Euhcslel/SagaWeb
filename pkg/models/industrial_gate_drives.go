package models

const TableNameIndustrialGateDrife = "industrial_gate_drives"

type IndustrialGateDrive struct {
	ID             int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name           string `gorm:"column:name;not null" json:"name"`
	UnitID         int32  `gorm:"column:unit_id" json:"unit_id"`
	WholesalePrice int32  `gorm:"column:wholesale_price" json:"wholesale_price"`
	RetailPrice    int32  `gorm:"column:retail_price" json:"retail_price"`
	Specifications string `gorm:"column:specifications" json:"specifications"`

	Unit Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*IndustrialGateDrive) TableName() string {
	return TableNameIndustrialGateDrife
}
