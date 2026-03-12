package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameResDrivesPriceHistory = "res_drives_price_history"

type ResDrivesPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	DriveID        int64           `gorm:"column:drive_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*ResDrivesPriceHistory) TableName() string {
	return TableNameResDrivesPriceHistory
}
