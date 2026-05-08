package types

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amount"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/sales_and_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/sales_and_products"
)

type SizeParams struct {
	MinValue int64
	MaxValue int64
}

type Config struct {
	CycleAmounts      []cycle_amount.CycleAmount
	LiftTypes         []lift_types.LiftType
	Colors            []colors.Color
	IndustrialDrives  []industrial_gate_drives.IndustrialGateDrive
	ResidentialDrives []residential_gate_drives.ResidentialGateDrive
	ManualDrive       manual_drive_prices.ManualDrivePrice
	WidthParams       SizeParams
	HeightParams      SizeParams
	Options           []options.Option
	Products          []products.Product
	DriveTypes        []enums.DriveType
	Rails             []rails.Rail
}

type Gate struct {
	Gate    sales_and_gates.SalesAndGate
	Options []options.Option
}

type Order struct {
	Gates    []Gate
	Products []sales_and_products.SalesAndProduct
	Status   enums.OrderStatus
}

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type UpdatedUserInfo struct {
	Fullname string
	Email    string
	Phone    string

	Company string
	Address string
}
