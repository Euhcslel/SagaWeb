package models

const TableNameSalesAndGate = "sales_and_gates"

type SalesAndGate struct {
	SaleID        int64 `gorm:"column:sale_id;primaryKey" json:"sale_id"`
	RowNumber     int64 `gorm:"column:row_number;primaryKey" json:"row_number"`
	GateTypeID      int32 `gorm:"column:gate_type;not null" json:"gate_type"`
	Width         int32 `gorm:"column:width;not null" json:"width"`
	Height        int32 `gorm:"column:height;not null" json:"height"`
	LiftTypeID    int64 `gorm:"column:lift_type_id;not null" json:"lift_type_id"`
	ColorInID     int64 `gorm:"column:color_in_id;not null" json:"color_in_id"`
	ColorOutID    int64 `gorm:"column:color_out_id;not null" json:"color_out_id"`
	CycleAmountID int64 `gorm:"column:cycle_amount_id;not null" json:"cycle_amount_id"`
	TotalPrice    int32 `gorm:"column:total_price;not null" json:"total_price"`
	StatusID      int32 `gorm:"column:status_id;not null" json:"status_id"`

	GateType GateType `gorm:"foreignKey:GateTypeID;references:ID"`
	LiftType LiftType `gorm:"foreignKey:LiftTypeID;references:ID"`
	ColorIn Color `gorm:"foreignKey:ColorInID;references:ID"`
	ColorOut Color `gorm:"foreignKey:ColorOutID;references:ID"`
	CycleAmount CycleAmount `gorm:"foreignKey:CycleAmountID;references:ID"`
	Status Status `gorm:"foreignKey:StatusID;references:ID"`
}

func (*SalesAndGate) TableName() string {
	return TableNameSalesAndGate
}
