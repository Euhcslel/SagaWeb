package repository

import (
	"project/internal/domain/gates_and_sales_manual_drive"
	"project/internal/domain/gates_and_sales_options"
	"project/internal/domain/industrial_gates_and_sales_drive"
	"project/internal/domain/managers_and_dealers"
	"project/internal/domain/residential_gates_and_sales_drive_rail"
	"project/internal/domain/sales"
	"project/internal/domain/sales_and_gates"
	"project/internal/domain/sales_and_products"
	"project/internal/domain/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAllDealerOrders(db *gorm.DB, user *users.User) ([]sales.Sale, error) {
	var orders []sales.Sale
	if err := db.
		Preload("Client").
		Preload("Manager").
		Where("client_id = ?", user.ID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func GetAllManagerOrders(db *gorm.DB, user *users.User) ([]sales.Sale, error) {
	var orders []sales.Sale
	if err := db.
		Preload("Client").
		Preload("Manager").
		Where("manager_id = ?", user.ID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrderGatesByOrderID(db *gorm.DB, saleID int64) ([]sales_and_gates.SalesAndGate, error) {
	var orderGates []sales_and_gates.SalesAndGate
	if err := db.
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("sale_id = ?", saleID).
		Find(&orderGates).Error; err != nil {
		return []sales_and_gates.SalesAndGate{}, err
	}

	return orderGates, nil
}

func GetAllOptionsForOrder(db *gorm.DB, gateIDs []int64) ([]gates_and_sales_options.GatesAndSalesOption, error) {
	var gateOptions []gates_and_sales_options.GatesAndSalesOption
	if err := db.
		Where("row_number IN ?", gateIDs).
		Preload("Option").
		Find(&gateOptions).Error; err != nil {
		return []gates_and_sales_options.GatesAndSalesOption{}, err
	}

	return gateOptions, nil
}

func GetAllOrderProducts(db *gorm.DB, saleID int64) ([]sales_and_products.SalesAndProduct, error) {
	var products []sales_and_products.SalesAndProduct
	if err := db.Preload("Product").Where("sale_id = ?", saleID).Find(&products).Error; err != nil {
		return []sales_and_products.SalesAndProduct{}, err
	}

	return products, nil
}

func GetCurrentGate(db *gorm.DB, saleID int64, gateID int64) (sales_and_gates.SalesAndGate, error) {
	var gate sales_and_gates.SalesAndGate
	if err := db.Model(sales_and_gates.SalesAndGate{}).
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("sale_id = ? and row_number = ?", saleID, gateID).
		First(&gate).Error; err != nil {
		return sales_and_gates.SalesAndGate{}, err
	}

	return gate, nil
}

func GetCurrentGateOptions(db *gorm.DB, saleID int64, gateID int64) ([]gates_and_sales_options.GatesAndSalesOption, error) {
	var options []gates_and_sales_options.GatesAndSalesOption
	if err := db.
		Model(gates_and_sales_options.GatesAndSalesOption{}).
		Preload("Option").Where("sale_id = ? and row_number = ?", saleID, gateID).
		Find(&options).Error; err != nil {
		return []gates_and_sales_options.GatesAndSalesOption{}, err
	}

	return options, nil
}

func GetIndustrialDriveForGate(db *gorm.DB, saleID int64, gateID int64) (*industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive, error) {
	var drive *industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive
	if err := db.
		Where("sale_id = ? AND row_number = ?", saleID, gateID).
		First(&drive).Error; err != nil {
		return &industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive{}, err
	}

	return drive, nil
}

func GetResidentialDriveForGate(db *gorm.DB, saleID int64, gateID int64) (*residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail, error) {
	var drive *residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail
	if err := db.
		Where("sale_id = ? AND row_number = ?", saleID, gateID).
		First(&drive).Error; err != nil {
		return &residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail{}, err
	}

	return drive, nil
}

func GetManualDriveForGate(db *gorm.DB, saleID int64, gateID int64) (*gates_and_sales_manual_drive.GatesAndSalesManualDrive, error) {
	var drive *gates_and_sales_manual_drive.GatesAndSalesManualDrive
	if err := db.
		Where("sale_id = ? AND row_number = ?", saleID, gateID).
		First(&drive).Error; err != nil {
		return &gates_and_sales_manual_drive.GatesAndSalesManualDrive{}, err
	}

	return drive, nil
}

func DeleteOrder(db *gorm.DB, sale sales.Sale) error {
	return db.Delete(&sale).Error
}

func CreateNewGate(db *gorm.DB, gate sales_and_gates.SalesAndGate) (sales_and_gates.SalesAndGate, error) {
	if err := db.Clauses(clause.Returning{}).Create(&gate).Error; err != nil {
		return sales_and_gates.SalesAndGate{}, err
	}

	return gate, nil
}

func CreateIndustrialDriveForGate(db *gorm.DB, saleID int64, rowNumber int64, driveID int32) error {
	if err := db.Create(&industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive{
		SaleID:    saleID,
		RowNumber: rowNumber,
		DriveID:   driveID,
	}).Error; err != nil {
		return err
	}

	return nil
}

func CreateResidentialDriveForGate(db *gorm.DB, saleID int64, rowNumber int64, driveID int32, railID int32) error {
	if err := db.Create(&residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail{
		SaleID:    saleID,
		RowNumber: rowNumber,
		DriveID:   driveID,
		RailID:    railID,
	}).Error; err != nil {
		return err
	}

	return nil
}

func CreateManualDriveForGate(db *gorm.DB, saleID int64, rowNumber int64, chainLength int32) error {
	if err := db.Create(&gates_and_sales_manual_drive.GatesAndSalesManualDrive{
		SaleID:      saleID,
		RowNumber:   rowNumber,
		ChainLength: chainLength,
	}).Error; err != nil {
		return err
	}

	return nil
}

func GetManagerIdByDealerId(db *gorm.DB, dealerID int64) (int64, error) {
	var managerID int64
	if err := db.Model(&managers_and_dealers.ManagerAndDealer{}).
		Select("manager_id").
		Where("dealer_id = ?", dealerID).
		Find(&managerID).Error; err != nil {
		return 0, err
	}

	return managerID, nil
}

func CreateNewOrder(db *gorm.DB, order *sales.Sale) error {
	return db.Create(&order).Error
}

func CreateOrderProducts(db *gorm.DB, products []sales_and_products.SalesAndProduct) error {
	return db.Create(&products).Error
}

func CreateGateOptions(db *gorm.DB, options []gates_and_sales_options.GatesAndSalesOption) error {
	return db.Create(&options).Error
}

func DeleteGateFromOrder(db *gorm.DB, gate sales_and_gates.SalesAndGate) error {
	return db.Delete(&gate).Error
}

func DeleteGateResidentialDrive(db *gorm.DB, saleID int64, gateID int64) error {
	var resGatesAndSalesDriveRail residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail
	if err := db.Where("sale_id = ? AND row_number = ?", saleID, gateID).First(&resGatesAndSalesDriveRail).Error; err != nil {
		return err
	}

	if err := db.Delete(&resGatesAndSalesDriveRail).Error; err != nil {
		return err
	}

	return nil
}

func DeleteGateIndustrialDrive(db *gorm.DB, saleID int64, gateID int64) error {
	var indGatesAndSalesDriveRail industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive
	if err := db.Where("sale_id = ? AND row_number = ?", saleID, gateID).First(&indGatesAndSalesDriveRail).Error; err != nil {
		return err
	}

	if err := db.Delete(&indGatesAndSalesDriveRail).Error; err != nil {
		return err
	}

	return nil
}

func DeleteGateManualDrive(db *gorm.DB, saleID int64, gateID int64) error {
	var gatesAndSalesManualDrive gates_and_sales_manual_drive.GatesAndSalesManualDrive
	if err := db.Where("sale_id = ? AND row_number = ?", saleID, gateID).First(&gatesAndSalesManualDrive).Error; err != nil {
		return err
	}

	if err := db.Delete(&gatesAndSalesManualDrive).Error; err != nil {
		return err
	}

	return nil
}

func UpdateGateIndustrialDrive(db *gorm.DB, saleID int64, gateID int64, driveID int32) error {
	var indGatesAndSalesDrive industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive
	if err := db.Where("sale_id = ? AND row_number = ?", saleID, gateID).First(&indGatesAndSalesDrive).Error; err != nil {
		return err
	}

	indGatesAndSalesDrive.DriveID = driveID

	if err := db.Save(&indGatesAndSalesDrive).Error; err != nil {
		return err
	}

	return nil
}

func UpdateGateResidentialDrive(db *gorm.DB, saleID int64, gateID int64, driveID int32, railID int32) error {
	var resGatesAndSalesDriveRail residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail
	if err := db.Where("sale_id = ? AND row_number = ?", saleID, gateID).First(&resGatesAndSalesDriveRail).Error; err != nil {
		return err
	}
	resGatesAndSalesDriveRail.DriveID = driveID
	resGatesAndSalesDriveRail.RailID = railID

	if err := db.Save(&resGatesAndSalesDriveRail).Error; err != nil {
		return err
	}

	return nil
}

func UpdateGateManualDrive(db *gorm.DB, saleID int64, gateID int64, chainLength int32) error {
	var gatesAndSalesManualDrive gates_and_sales_manual_drive.GatesAndSalesManualDrive
	if err := db.Where("sale_id = ? AND row_number = ?", saleID, gateID).First(&gatesAndSalesManualDrive).Error; err != nil {
		return err
	}
	gatesAndSalesManualDrive.ChainLength = chainLength
	if err := db.Save(&gatesAndSalesManualDrive).Error; err != nil {
		return err
	}

	return nil
}

func UpdateGate(db *gorm.DB, gate sales_and_gates.SalesAndGate) error {
	return db.Save(&gate).Error
}
