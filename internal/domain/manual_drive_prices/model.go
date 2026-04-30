package manual_drive_prices

import "github.com/shopspring/decimal"

const TableNameManualDrivePrices = "manual_drive_prices"

type ManualDrivePrice struct {
	ID int64 `gorm:"column:id;primaryKey"`

	ChainMeterRetailPrice    decimal.Decimal `gorm:"column:chain_meter_retail_price"`
	ChainMeterWholesalePrice decimal.Decimal `gorm:"column:chain_meter_wholesale_price"`

	RcpRetailPrice    decimal.Decimal `gorm:"column:rcp_retail_price"`
	RcpWholesalePrice decimal.Decimal `gorm:"column:rcp_wholesale_price"`
}

func (*ManualDrivePrice) TableName() string {
	return TableNameManualDrivePrices
}
