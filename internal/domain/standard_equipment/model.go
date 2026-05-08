package standard_equipment

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
)

const TableNameStandardEquipment = "standard_equipment"

type StandardEquipment struct {
	GateType enums.GateType `gorm:"column:gate_type;primaryKey"`
	OptionID int64          `gorm:"column:option_id;primaryKey"`
	Amount   int32          `gorm:"column:amount"`

	Option options.Option `gorm:"foreignKey:OptionID;references:ID"`
}

func (*StandardEquipment) TableName() string {
	return TableNameStandardEquipment
}
