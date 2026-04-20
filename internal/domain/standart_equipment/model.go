package standart_equipment

import (
	"project/internal/domain/options"
)

const TableNameStandartEquipment = "standart_equipment"

type StandartEquipment struct {
	GateTypeID int64 `gorm:"column:gate_type_id;primaryKey"`
	OptionID   int64 `gorm:"column:option_id;primaryKey"`
	Amount     int32 `gorm:"column:amount"`

	Option options.Option `gorm:"foreignKey:OptionID;references:ID"`
}

func (*StandartEquipment) TableName() string {
	return TableNameStandartEquipment
}
