package ind_drives_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameIndDrivesPriceHistory = "ind_drives_price_history"

type IndDrivesPriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	DriveID        int64           `gorm:"column:drive_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*IndDrivesPriceHistory) TableName() string {
	return TableNameIndDrivesPriceHistory
}
