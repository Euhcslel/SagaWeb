// Package repository содержит функции для работы с базой данных

package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
)

// CreateNewColor сохраняет новый цвет в базе данных.
func CreateNewColor(color colors.Color) error {
	return database.DB.Create(&color).Error
}

// CreateNewCycleAmount сохраняет новое количество циклов в базе данных.
func CreateNewCycleAmount(cycleAmount cycle_amounts.CycleAmount) error {
	return database.DB.Create(&cycleAmount).Error
}

// CreateNewIndustrialDrive сохраняет новый промышленный привод в базе данных.
func CreateNewIndustrialDrive(industrialDrive industrial_gate_drives.IndustrialGateDrive) error {
	return database.DB.Create(&industrialDrive).Error
}

// CreateNewLiftType сохраняет новый тип подъема в базе данных.
func CreateNewLiftType(liftType lift_types.LiftType) error {
	return database.DB.Create(&liftType).Error
}

// CreateNewOption сохраняет новую дополнительную опцию в базе данных.
func CreateNewOption(option options.Option) error {
	return database.DB.Create(&option).Error
}

// CreateNewProduct сохраняет новый товар в базе данных.
func CreateNewProduct(product products.Product) error {
	return database.DB.Create(&product).Error
}

// CreateNewRail сохраняет новую направляющую в базе данных.
func CreateNewRail(rail rails.Rail) error {
	return database.DB.Create(&rail).Error
}

// CreateNewResidentialDrive сохраняет новый бытовой привод в базе данных.
func CreateNewResidentialDrive(residentialDrive residential_gate_drives.ResidentialGateDrive) error {
	return database.DB.Create(&residentialDrive).Error
}

// UpdateColor обновляет информацию о цвете в базе данных.
func UpdateColor(color colors.Color) error {
	return database.DB.Model(&colors.Color{}).
		Where("id = ?", color.ID).
		Updates(map[string]any{
			"code": color.Code,
			"hex":  color.Hex,
		}).Error
}

// UpdateCycleAmount обновляет информацию о количестве циклов в базе данных.
func UpdateCycleAmount(cycleAmount cycle_amounts.CycleAmount) error {
	return database.DB.Model(&cycle_amounts.CycleAmount{}).
		Where("id = ?", cycleAmount.ID).
		Updates(map[string]any{
			"amount":           cycleAmount.Amount,
			"wholesale_markup": cycleAmount.WholesaleMarkup,
			"retail_markup":    cycleAmount.RetailMarkup,
		}).Error
}

// UpdateIndustrialDrive обновляет информацию о промышленном приводе в базе данных.
func UpdateIndustrialDrive(industrialDrive industrial_gate_drives.IndustrialGateDrive) error {
	return database.DB.Model(&industrial_gate_drives.IndustrialGateDrive{}).
		Where("id = ?", industrialDrive.ID).
		Updates(map[string]any{
			"name":            industrialDrive.Name,
			"wholesale_price": industrialDrive.WholesalePrice,
			"retail_price":    industrialDrive.RetailPrice,
			"specifications":  industrialDrive.Specifications,
			"image_path":      industrialDrive.ImagePath,
		}).Error
}

// UpdateLiftType обновляет информацию о типе подъема в базе данных.
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

// UpdateOption обновляет информацию о дополнительной опции в базе данных.
func UpdateOption(option options.Option) error {
	return database.DB.Model(&options.Option{}).
		Where("id = ?", option.ID).
		Updates(map[string]any{
			"name":            option.Name,
			"wholesale_price": option.WholesalePrice,
			"retail_price":    option.RetailPrice,
			"image_path":      option.ImagePath,
		}).Error
}

// UpdateProduct обновляет информацию о товаре в базе данных.
func UpdateProduct(product products.Product) error {
	return database.DB.Model(&products.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]any{
			"name":            product.Name,
			"wholesale_price": product.WholesalePrice,
			"retail_price":    product.RetailPrice,
			"image_path":      product.ImagePath,
		}).Error
}

// UpdateRail обновляет информацию о направляющей в базе данных.
func UpdateRail(rail rails.Rail) error {
	return database.DB.Model(&rails.Rail{}).
		Where("id = ?", rail.ID).
		Updates(map[string]any{
			"name":            rail.Name,
			"wholesale_price": rail.WholesalePrice,
			"retail_price":    rail.RetailPrice,
			"specifications":  rail.Specifications,
			"image_path":      rail.ImagePath,
		}).Error
}

// UpdateResidentialDrive обновляет информацию о бытовом приводе в базе данных.
func UpdateResidentialDrive(residentialDrive residential_gate_drives.ResidentialGateDrive) error {
	return database.DB.Model(&residential_gate_drives.ResidentialGateDrive{}).
		Where("id = ?", residentialDrive.ID).
		Updates(map[string]any{
			"name":            residentialDrive.Name,
			"wholesale_price": residentialDrive.WholesalePrice,
			"retail_price":    residentialDrive.RetailPrice,
			"specifications":  residentialDrive.Specifications,
			"image_path":      residentialDrive.ImagePath,
		}).Error
}

// DeleteColor удаляет цвет из базы данных по его идентификатору.
func DeleteColor(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&colors.Color{}).Error
}

// DeleteCycleAmount удаляет информацию о количестве циклов из базы данных по его идентификатору.
func DeleteCycleAmount(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&cycle_amounts.CycleAmount{}).Error
}

// DeleteIndustrialDrive удаляет информацию о промышленном приводе из базы данных по его идентификатору.
func DeleteIndustrialDrive(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&industrial_gate_drives.IndustrialGateDrive{}).Error
}

// DeleteLiftType удаляет информацию о типе подъема из базы данных по его идентификатору.
func DeleteLiftType(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&lift_types.LiftType{}).Error
}

// DeleteOption удаляет информацию о дополнительной опции из базы данных по её идентификатору.
func DeleteOption(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&options.Option{}).Error
}

// DeleteProduct удаляет информацию о товаре из базы данных по его идентификатору.
func DeleteProduct(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&products.Product{}).Error
}

// DeleteRail удаляет информацию о направляющей из базы данных по её идентификатору.
func DeleteRail(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&rails.Rail{}).Error
}

// DeleteResidentialDrive удаляет информацию о бытовом приводе из базы данных по его идентификатору.
func DeleteResidentialDrive(rowId int64) error {
	return database.DB.Where("id = ?", rowId).Delete(&residential_gate_drives.ResidentialGateDrive{}).Error
}

// GetAllUsers возвращает список всех пользователей из базы данных.
func GetAllUsers() ([]users.User, error) {
	var result []users.User
	err := database.DB.Find(&result).Error
	return result, err
}

// AdminCreateUser создает нового пользователя в базе данных.
func AdminCreateUser(user users.User) error {
	return database.DB.Create(&user).Error
}

// AdminUpdateUser обновляет информацию о пользователе в базе данных.
func AdminUpdateUser(user users.User) error {
	return database.DB.Model(&users.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"fullname":     user.Fullname,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"role":         user.Role,
		}).Error
}

// AdminUpdateUserPassword обновляет хэш пароля пользователя в базе данных.
func AdminUpdateUserPassword(id int64, hash []byte) error {
	return database.DB.Model(&users.User{}).
		Where("id = ?", id).
		Update("password_hash", hash).Error
}

// AdminDeleteUser удаляет пользователя из базы данных по его идентификатору.
func AdminDeleteUser(id int64) error {
	return database.DB.Where("id = ?", id).Delete(&users.User{}).Error
}

// GetAllManualDrivePrices возвращает список всех цен на ручные приводы из базы данных.
func GetAllManualDrivePrices() ([]manual_drive_prices.ManualDrivePrice, error) {
	var result []manual_drive_prices.ManualDrivePrice
	if err := database.DB.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// CreateManualDrivePrice создает новую цену на ручной привод в базе данных.
func CreateManualDrivePrice(price manual_drive_prices.ManualDrivePrice) error {
	return database.DB.Create(&price).Error
}

// UpdateManualDrivePrice обновляет информацию о цене на ручной привод в базе данных.
func UpdateManualDrivePrice(price manual_drive_prices.ManualDrivePrice) error {
	return database.DB.Model(&manual_drive_prices.ManualDrivePrice{}).
		Where("id = ?", price.ID).
		Updates(map[string]any{
			"chain_meter_retail_price":    price.ChainMeterRetailPrice,
			"chain_meter_wholesale_price": price.ChainMeterWholesalePrice,
			"rcp_retail_price":            price.RcpRetailPrice,
			"rcp_wholesale_price":         price.RcpWholesalePrice,
		}).Error
}

// DeleteManualDrivePrice удаляет информацию о цене на ручной привод из базы данных по его идентификатору.
func DeleteManualDrivePrice(id int64) error {
	return database.DB.Where("id = ?", id).Delete(&manual_drive_prices.ManualDrivePrice{}).Error
}
