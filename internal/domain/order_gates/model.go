// Package order_gates предоставляет модель для работы с воротами в заказе.
// Данные о воротах в заказе хранятся в таблице order_gates.

package order_gates

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/orders"

	"github.com/shopspring/decimal"
)

const TableNameOrderGate = "order_gates"

type OrderGate struct {
	OrderID            int64           `gorm:"column:order_id;primaryKey"`
	RowNumber          int64           `gorm:"column:row_number;primaryKey"`
	GateType           enums.GateType  `gorm:"type:gate_type;not null;"`
	Width              int32           `gorm:"column:width;not null"`
	Height             int32           `gorm:"column:height;not null"`
	Headroom           int32           `gorm:"column:headroom;not null"`
	LiftTypeID         int64           `gorm:"column:lift_type_id;not null"`
	ColorOutID         int64           `gorm:"column:color_out_id;not null"`
	CycleAmountID      int64           `gorm:"column:cycle_amount_id;not null"`
	DriveType          enums.DriveType `gorm:"column:drive_type;not null"`
	Amount             int32           `gorm:"column:amount;not null"`
	GateRetailPrice    decimal.Decimal `gorm:"column:gate_retail_price;type:decimal(10,2);not null"`
	GateWholesalePrice decimal.Decimal `gorm:"column:gate_wholesale_price;type:decimal(10,2);not null"`

	LiftType    lift_types.LiftType       `gorm:"foreignKey:LiftTypeID;references:ID"`
	ColorOut    colors.Color              `gorm:"foreignKey:ColorOutID;references:ID"`
	CycleAmount cycle_amounts.CycleAmount `gorm:"foreignKey:CycleAmountID;references:ID"`
	Order       orders.Order              `gorm:"foreignKey:OrderID;references:ID"`
}

func (*OrderGate) TableName() string {
	return TableNameOrderGate
}
