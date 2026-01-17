package models

const TableNameProduct = "products"

type Product struct {
	ID             int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name           string `gorm:"column:name;not null" json:"name"`
	UnitID         int32  `gorm:"column:unit_id" json:"unit_id"`
	WholesalePrice int32  `gorm:"column:wholesale_price" json:"wholesale_price"`
	RetailPrice    int32  `gorm:"column:retail_price" json:"retail_price"`

	Unit Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*Product) TableName() string {
	return TableNameProduct
}
