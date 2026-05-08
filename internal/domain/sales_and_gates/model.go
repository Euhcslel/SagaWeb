package sales_and_gates

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amount"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/sales"

	"github.com/shopspring/decimal"
)

const TableNameSalesAndGate = "sales_and_gates"

type SalesAndGate struct {
	SaleID             int64           `gorm:"column:sale_id;primaryKey"`
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

	LiftType    lift_types.LiftType      `gorm:"foreignKey:LiftTypeID;references:ID"`
	ColorOut    colors.Color             `gorm:"foreignKey:ColorOutID;references:ID"`
	CycleAmount cycle_amount.CycleAmount `gorm:"foreignKey:CycleAmountID;references:ID"`
	Sale        sales.Sale               `gorm:"foreignKey:SaleID;references:ID"`
}

func (*SalesAndGate) TableName() string {
	return TableNameSalesAndGate
}
