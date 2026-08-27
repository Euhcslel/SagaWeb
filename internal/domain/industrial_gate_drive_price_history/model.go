// Package industrial_gate_drive_price_history предоставляет модель для работы с историей цен за промышленные привода.
// История цен за промышленные привода хранится в таблице industrial_gate_drive_price_history.

package industrial_gate_drive_price_history

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameIndustrialGateDrivePriceHistory = "industrial_gate_drive_price_history"

type IndustrialGateDrivePriceHistory struct {
	ID             int64           `gorm:"column:id;primaryKey"`
	DriveID        int64           `gorm:"column:drive_id"`
	WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
	SetAt          time.Time       `gorm:"column:set_at"`
}

func (*IndustrialGateDrivePriceHistory) TableName() string {
	return TableNameIndustrialGateDrivePriceHistory
}
