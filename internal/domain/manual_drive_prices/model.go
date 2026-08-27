// Package manual_drive_prices предоставляет модель для работы с ценами за ручной тип подъема.
// Данные о ценах за ручной тип подъема хранятся в таблице manual_drive_prices.

package manual_drive_prices

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameManualDrivePrices = "manual_drive_prices"

type ManualDrivePrice struct {
	ID int64 `gorm:"column:id;primaryKey"`

	ChainMeterRetailPrice    decimal.Decimal `gorm:"column:chain_meter_retail_price"`
	ChainMeterWholesalePrice decimal.Decimal `gorm:"column:chain_meter_wholesale_price"`

	RcpRetailPrice    decimal.Decimal `gorm:"column:rcp_retail_price"`
	RcpWholesalePrice decimal.Decimal `gorm:"column:rcp_wholesale_price"`

	CreatedAt time.Time `gorm:"column:created_at;default:NOW();not null"`
}

func (*ManualDrivePrice) TableName() string {
	return TableNameManualDrivePrices
}
