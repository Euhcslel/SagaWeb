package types

import (
	"project/internal/domain/colors"
	"project/internal/domain/cycle_amount"
	"project/internal/domain/enums"
	"project/internal/domain/industrial_gate_drives"
	"project/internal/domain/lift_types"
	"project/internal/domain/options"
	"project/internal/domain/products"
	"project/internal/domain/rails"
	"project/internal/domain/residential_gate_drives"
	"project/internal/domain/sales_and_gates"
	"project/internal/domain/sales_and_products"
)

type SizeParams struct {
	MinValue uint64
	MaxValue uint64
}

type Config struct {
	CycleAmounts      []cycle_amount.CycleAmount
	LiftTypes         []lift_types.LiftType
	Colors            []colors.Color
	IndustrialDrives  []industrial_gate_drives.IndustrialGateDrive
	ResidentialDrives []residential_gate_drives.ResidentialGateDrive
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
}

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)
