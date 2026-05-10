package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/gate_sizes"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/types"
)

func GetSizeForDimensions(width, height int64, gateType enums.GateType) (*gate_sizes.GateSize, error) {
	var size gate_sizes.GateSize

	// Ищем наименьший стандартный размер, в который вписываются запрошенные габариты.
	// Если точного совпадения нет — берём ближайший больший.
	if err := database.DB.Model(&gate_sizes.GateSize{}).
		Where("width >= ? AND height >= ?", width, height).
		Where("gate_type = ?", gateType).
		Limit(1).
		Order("width asc, height asc").
		First(&size).Error; err != nil {
		return nil, err
	}
	return &size, nil
}

func GetIndustrialDrives() ([]industrial_gate_drives.IndustrialGateDrive, error) {
	var drives []industrial_gate_drives.IndustrialGateDrive
	if err := database.DB.Find(&drives).Error; err != nil {
		return nil, err
	}
	return drives, nil
}

func GetResidentialDrives() ([]residential_gate_drives.ResidentialGateDrive, error) {
	var drives []residential_gate_drives.ResidentialGateDrive
	if err := database.DB.Find(&drives).Error; err != nil {
		return nil, err
	}
	return drives, nil
}

func GetRails() ([]rails.Rail, error) {
	var rails []rails.Rail
	if err := database.DB.Find(&rails).Error; err != nil {
		return nil, err
	}
	return rails, nil
}

func GetLiftTypes() ([]lift_types.LiftType, error) {
	var liftTypes []lift_types.LiftType
	if err := database.DB.Find(&liftTypes).Error; err != nil {
		return nil, err
	}
	return liftTypes, nil
}

func GetColors() ([]colors.Color, error) {
	var colors []colors.Color
	if err := database.DB.Find(&colors).Error; err != nil {
		return nil, err
	}
	return colors, nil
}

func GetCycleAmounts() ([]cycle_amounts.CycleAmount, error) {
	var cycleAmounts []cycle_amounts.CycleAmount
	if err := database.DB.Find(&cycleAmounts).Error; err != nil {
		return nil, err
	}
	return cycleAmounts, nil
}

func GetMaxAndMinWidth(gateType enums.GateType) (*types.SizeParams, error) {
	var widthParams types.SizeParams
	if err := database.DB.
		Model(&gate_sizes.GateSize{}).
		Select("MAX(width) as max_value, MIN(width) as min_value").
		Where("gate_type = ?", gateType).
		Scan(&widthParams).Error; err != nil {
		return nil, err
	}
	return &widthParams, nil
}

func GetMaxAndMinHeight(gateType enums.GateType) (*types.SizeParams, error) {
	var heightParams types.SizeParams
	if err := database.DB.
		Model(&gate_sizes.GateSize{}).
		Select("MAX(height) as max_value, MIN(height) as min_value").
		Where("gate_type = ?", gateType).
		Scan(&heightParams).Error; err != nil {
		return nil, err
	}
	return &heightParams, nil
}

func GetProducts() ([]products.Product, error) {
	var products []products.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func GetOptions() ([]options.Option, error) {
	var options []options.Option
	if err := database.DB.Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func GetManualDrivePrices() (*manual_drive_prices.ManualDrivePrice, error) {
	var manualDrivePrice manual_drive_prices.ManualDrivePrice
	if err := database.DB.First(&manualDrivePrice).Error; err != nil {
		return nil, err
	}
	return &manualDrivePrice, nil
}
