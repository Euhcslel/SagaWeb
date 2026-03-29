package models

const TableNameSalesAndGate = "sales_and_gates"

type SalesAndGate struct {
	SaleID        int64     `gorm:"column:sale_id;primaryKey"`
	RowNumber     int64     `gorm:"column:row_number;primaryKey"`
	GateType      GateType  `gorm:"type:gate_type;not null;"`
	Width         int32     `gorm:"column:width;not null"`
	Height        int32     `gorm:"column:height;not null"`
	Headroom      int32     `gorm:"column:headroom;not null"`
	LiftTypeID    int64     `gorm:"column:lift_type_id;not null"`
	ColorOutID    int64     `gorm:"column:color_out_id;not null"`
	CycleAmountID int64     `gorm:"column:cycle_amount_id;not null"`
	DriveType     DriveType `gorm:"column:drive_type;not null"`
	Amount        int32     `gorm:"column:amount;not null"`

	LiftType    LiftType    `gorm:"foreignKey:LiftTypeID;references:ID"`
	ColorOut    Color       `gorm:"foreignKey:ColorOutID;references:ID"`
	CycleAmount CycleAmount `gorm:"foreignKey:CycleAmountID;references:ID"`
}

func (*SalesAndGate) TableName() string {
	return TableNameSalesAndGate
}
