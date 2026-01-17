package models

const TableNameOption = "options"

type Option struct {
	ID             int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name           string `gorm:"column:name;not null" json:"name"`
	UnitID         int32  `gorm:"column:unit_id" json:"unit_id"`
	WholesalePrice int32  `gorm:"column:wholesale_price" json:"wholesale_price"`
	RetailPrice    int32  `gorm:"column:retail_price" json:"retail_price"`
	ForSale        bool   `gorm:"column:for_sale;not null" json:"for_sale"`
	Condition      string `gorm:"column:condition" json:"condition"`

	Unit Unit `gorm:"foreignKey:UnitID;references:ID"`
}

func (*Option) TableName() string {
	return TableNameOption
}
