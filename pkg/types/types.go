package types

import "project/pkg/models"

type WidthParams struct {
	MinWidth uint64
	MaxWidth uint64
}

type HeightParams struct {
	MinHeight uint64
	MaxHeight uint64
}

type Config struct {
	CycleAmounts []models.CycleAmount
	LiftTypes    []models.LiftType
	Colors       []models.Color
	MontageTypes []models.MontageType
	IndustrialDrives []models.IndustrialGateDrive
	ResidentialDrives []models.ResidentialGateDrive
	WidthParams  WidthParams
	HeightParams HeightParams
}

type Gate struct {
	Gate models.SalesAndGate
	Options []models.Option
}

type Order struct {
	Gates []Gate
	Products []models.Product
}

type GateConfig struct {
    GateType   string  `json:"gateType"`
    Width      int     `json:"width,string"`
    Height     int     `json:"height,string"`
    LiftType   string  `json:"liftType"`
    ColorIn    string  `json:"colorIn"`
    ColorOut   string  `json:"colorOut"`
    Drive      string  `json:"drive"`
    CycleAmount string `json:"cycleAmount"`
    GatePrice  float64 `json:"gatePrice"`
}

type Products map[string]any

type OrderRequest struct {
    OrderGates []GateConfig `json:"orderGates"`
    Products   Products     `json:"products"`
}
