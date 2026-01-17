package models

const TableNameStandartEquipment = "standart_equipment"

type StandartEquipment struct {
	GateTypeID int32 `gorm:"column:gate_type_id;primaryKey" json:"gate_type_id"`
	OptionID   int32 `gorm:"column:option_id;primaryKey" json:"option_id"`
	Amount     int32 `gorm:"column:amount" json:"amount"`

	Option Option `gorm:"foreignKey:OptionID;references:ID"`
}

func (*StandartEquipment) TableName() string {
	return TableNameStandartEquipment
}
