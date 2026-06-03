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

func UpdateColor(color colors.Color) error {
	return database.DB.Model(&colors.Color{}).
		Where("id = ?", color.ID).
		Updates(map[string]any{
			"code": color.Code,
			"hex":  color.Hex,
		}).Error
}

func UpdateCycleAmount(cycleAmount cycle_amounts.CycleAmount) error {
	return database.DB.Model(&cycle_amounts.CycleAmount{}).
		Where("id = ?", cycleAmount.ID).
		Updates(map[string]any{
			"amount":           cycleAmount.Amount,
			"wholesale_markup": cycleAmount.WholesaleMarkup,
			"retail_markup":    cycleAmount.RetailMarkup,
		}).Error
}

func UpdateIndustrialDrive(industrialDrive industrial_gate_drives.IndustrialGateDrive) error {
	return database.DB.Model(&industrial_gate_drives.IndustrialGateDrive{}).
		Where("id = ?", industrialDrive.ID).
		Updates(map[string]any{
			"name":             industrialDrive.Name,
			"wholesale_price":  industrialDrive.WholesalePrice,
			"retail_price":     industrialDrive.RetailPrice,
			"specifications":   industrialDrive.Specifications,
		}).Error
}

func UpdateLiftType(liftType lift_types.LiftType) error {
	return database.DB.Model(&lift_types.LiftType{}).
		Where("id = ?", liftType.ID).
		Updates(map[string]any{
			"name":             liftType.Name,
			"min_headroom":     liftType.MinHeadroom,
			"max_headroom":     liftType.MaxHeadroom,
			"wholesale_markup": liftType.WholesaleMarkup,
			"retail_markup":    liftType.RetailMarkup,
		}).Error
}

func UpdateOption(option options.Option) error {
	return database.DB.Model(&options.Option{}).
		Where("id = ?", option.ID).
		Updates(map[string]any{
			"name":            option.Name,
			"wholesale_price": option.WholesalePrice,
			"retail_price":    option.RetailPrice,
		}).Error
}

func UpdateProduct(product products.Product) error {
	return database.DB.Model(&products.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]any{
			"name":            product.Name,
			"wholesale_price": product.WholesalePrice,
			"retail_price":    product.RetailPrice,
		}).Error
}

func UpdateResidentialDrive(residentialDrive residential_gate_drives.ResidentialGateDrive) error {
	return database.DB.Model(&residential_gate_drives.ResidentialGateDrive{}).
		Where("id = ?", residentialDrive.ID).
		Updates(map[string]any{
			"name":            residentialDrive.Name,
			"wholesale_price": residentialDrive.WholesalePrice,
			"retail_price":    residentialDrive.RetailPrice,
			"specifications":  residentialDrive.Specifications,
		}).Error
}

func DeleteColor(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&colors.Color{}).Error
}

func DeleteCycleAmount(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&cycle_amounts.CycleAmount{}).Error
}

func DeleteIndustrialDrive(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&industrial_gate_drives.IndustrialGateDrive{}).Error
}

func DeleteLiftType(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&lift_types.LiftType{}).Error
}

func DeleteOption(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&options.Option{}).Error
}

func DeleteProduct(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&products.Product{}).Error
}

func DeleteResidentialDrive(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&residential_gate_drives.ResidentialGateDrive{}).Error
}
