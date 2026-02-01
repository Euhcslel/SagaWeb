package models

const TableNameResidentialGatesAndSalesDriveRail = "residential_gates_and_sales_drive_rail"

type ResidentialGatesAndSalesDriveRail struct {
	SaleID    int64 `gorm:"column:sale_id;primaryKey"`
	RowNumber int64 `gorm:"column:row_number;primaryKey"`
	DriveID   int32 `gorm:"column:drive_id;not null"`
	RailID    int32 `gorm:"column:rail_id;not null"`

	Drive ResidentialGateDrive `gorm:"foreignKey:DriveID;references:ID"`
	Rail  Rail                 `gorm:"foreignKey:RailID;references:ID"`
}

func (*ResidentialGatesAndSalesDriveRail) TableName() string {
	return TableNameResidentialGatesAndSalesDriveRail
}
