// Package types содержит структуры данных, используемые в приложении.

package types

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_products"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
)

// PricePair представляет пару цен: розничную и оптовую
type PricePair struct {
	RetailPrice    decimal.Decimal
	WholesalePrice decimal.Decimal
}

// SizeParams представляет параметры размера
type SizeParams struct {
	MinValue int64
	MaxValue int64
}

// Config представляет конфигурацию ворот для калькулятора.
// Содержит все необходимые данные для расчета стоимости ворот.
type Config struct {
	CycleAmounts      []cycle_amounts.CycleAmount
	LiftTypes         []lift_types.LiftType
	Colors            []colors.Color
	IndustrialDrives  []industrial_gate_drives.IndustrialGateDrive
	ResidentialDrives []residential_gate_drives.ResidentialGateDrive
	ManualDrive       *manual_drive_prices.ManualDrivePrice
	WidthParams       *SizeParams
	HeightParams      *SizeParams
	Options           []options.Option
	Products          []products.Product
	DriveTypes        []enums.DriveType
	Rails             []rails.Rail
}

// Gate представляет ворота с приводом и опциями.
type Gate struct {
	Gate    order_gates.OrderGate
	Drive   any
	Options []options.Option
}

// Order представляет заказ с воротами, товарами и статусом.
type Order struct {
	Gates           []Gate
	Products        []order_products.OrderProduct
	Status          enums.OrderStatus
	ManufactureDate *time.Time
}

// UpdatedUserInfo представляет обновленную информацию о пользователе.
type UpdatedUserInfo struct {
	Fullname string
	Email    string
	Phone    string

	Company string
	Address string

	NewPassword  string
	PasswordHash []byte
}

// ResidentialDriveRail представляет комбинацию привода и направляющей.
type ResidentialDriveRail struct {
	Drive residential_gate_drives.ResidentialGateDrive
	Rail  rails.Rail
}

// ManualDrive представляет ручной привод с длиной цепи и информацией о цене.
type ManualDrive struct {
	ChainLength int32
	PriceInfo   manual_drive_prices.ManualDrivePrice
}
