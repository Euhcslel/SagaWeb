package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
)

func CreateNewColor(color colors.Color) error {
	return database.DB.Create(&color).Error
}

func CreateNewCycleAmount(cycleAmount cycle_amounts.CycleAmount) error {
	return database.DB.Create(&cycleAmount).Error
}

func CreateNewIndustrialDrive(industrialDrive industrial_gate_drives.IndustrialGateDrive) error {
	return database.DB.Create(&industrialDrive).Error
}

func CreateNewLiftType(liftType lift_types.LiftType) error {
	return database.DB.Create(&liftType).Error
}

func CreateNewOption(option options.Option) error {
	return database.DB.Create(&option).Error
}

func CreateNewProduct(product products.Product) error {
	return database.DB.Create(&product).Error
}

func CreateNewResidentialDrive(residentialDrive residential_gate_drives.ResidentialGateDrive) error {
	return database.DB.Create(&residentialDrive).Error
}
