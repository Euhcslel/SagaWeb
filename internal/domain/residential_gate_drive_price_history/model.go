package residential_gate_drive_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameResidentialGateDrivePriceHistory = "residential_gate_drive_price_history"

type ResidentialGateDrivePriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	DriveID        int64           `gorm:"column:drive_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*ResidentialGateDrivePriceHistory) TableName() string {
	return TableNameResidentialGateDrivePriceHistory
}
