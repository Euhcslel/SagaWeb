package types

import "project/pkg/models"

type SizeParams struct {
	MinValue uint64
	MaxValue uint64
}

type Config struct {
	CycleAmounts      []models.CycleAmount
	LiftTypes         []models.LiftType
	Colors            []models.Color
	IndustrialDrives  []models.IndustrialGateDrive
	ResidentialDrives []models.ResidentialGateDrive
	WidthParams       SizeParams
	HeightParams      SizeParams
	Options           []models.Option
	Products          []models.Product
	DriveTypes        []models.DriveType
	Rails             []models.Rail
}

type Gate struct {
	Gate    models.SalesAndGate
	Options []models.Option
}

type Order struct {
	Gates    []Gate
	Products []models.SalesAndProduct
}
