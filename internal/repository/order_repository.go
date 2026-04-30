package repository

import (
	"project/internal/database"
	"project/internal/domain/cycle_amount"
	"project/internal/domain/enums"
	"project/internal/domain/gates_and_sales_manual_drive"
	"project/internal/domain/gates_and_sales_options"
	"project/internal/domain/industrial_gate_drives"
	"project/internal/domain/industrial_gates_and_sales_drive"
	"project/internal/domain/lift_types"
	"project/internal/domain/managers_and_dealers"
	"project/internal/domain/options"
	"project/internal/domain/products"
	"project/internal/domain/rails"
	"project/internal/domain/residential_gate_drives"
	"project/internal/domain/residential_gates_and_sales_drive_rail"
	"project/internal/domain/sales"
	"project/internal/domain/sales_and_bills"
	"project/internal/domain/sales_and_contracts"
	"project/internal/domain/sales_and_documents"
	"project/internal/domain/sales_and_gates"
	"project/internal/domain/sales_and_offers"
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

func CreateIndustrialDriveForGate(db *gorm.DB, saleID int64, rowNumber int64, driveID int64) error {
	if err := db.Create(&industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive{
		SaleID:    saleID,
		RowNumber: rowNumber,
		DriveID:   driveID,
	}).Error; err != nil {
		return err
	}

	return nil
}

func CreateResidentialDriveForGate(db *gorm.DB, saleID int64, rowNumber int64, driveID int64, railID int64) error {
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

func UpdateGateIndustrialDrive(db *gorm.DB, saleID int64, gateID int64, driveID int64) error {
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

func UpdateGateResidentialDrive(db *gorm.DB, saleID int64, gateID int64, driveID int64, railID int64) error {
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

func UpdateOrderStatus(db *gorm.DB, sale *sales.Sale, status enums.OrderStatus) error {
	sale.Status = status
	return db.Save(&sale).Error
}

func GetDocumentsNameList(db *gorm.DB, saleID int64) ([]string, error) {
	var documentsNameList []string
	if err := db.Model(&sales_and_documents.SalesAndDocument{}).
		Where("sale_id = ?", saleID).
		Pluck("name", &documentsNameList).Error; err != nil {
		return nil, err
	}

	return documentsNameList, nil
}

func GetOffersNumberList(db *gorm.DB, saleID int64) ([]string, error) {
	var offersNumberList []string
	if err := db.Model(&sales_and_offers.SalesAndOffer{}).
		Where("sale_id = ?", saleID).
		Pluck("offer_number", &offersNumberList).Error; err != nil {
		return nil, err
	}

	return offersNumberList, nil
}

func GetContractsNumberList(db *gorm.DB, saleID int64) ([]string, error) {
	var contractsNumberList []string
	if err := db.Model(&sales_and_contracts.SalesAndContract{}).
		Where("sale_id = ?", saleID).
		Pluck("contract_number", &contractsNumberList).Error; err != nil {
		return nil, err
	}

	return contractsNumberList, nil
}

func GetBillsNumberList(db *gorm.DB, saleID int64) ([]string, error) {
	var billsNumberList []string
	if err := db.Model(&sales_and_bills.SalesAndBill{}).
		Where("sale_id = ?", saleID).
		Pluck("bill_number", &billsNumberList).Error; err != nil {
		return nil, err
	}

	return billsNumberList, nil
}

func GetOrderStatus(user *users.User, saleID int64) (string, error) {
	var status string
	if err := database.DB.Model(&sales.Sale{}).
		Where("id = ?", saleID).
		Pluck("status", &status).Error; err != nil {
		return "", err
	}
	return status, nil
}

func AttachOfferToOrder(db *gorm.DB, saleID int64, offerNumber string, path string) error {
	offerAndSale := sales_and_offers.SalesAndOffer{
		SaleID:      saleID,
		OfferNumber: offerNumber,
		Path:        path,
	}
	return db.Create(&offerAndSale).Error
}

func AttachContractToOrder(db *gorm.DB, saleID int64, contractNumber string, path string) error {
	contractAndSale := sales_and_contracts.SalesAndContract{
		SaleID:         saleID,
		ContractNumber: contractNumber,
		Path:           path,
	}
	return db.Create(&contractAndSale).Error
}

func AttachBillToOrder(db *gorm.DB, saleID int64, billNumber string, path string) error {
	billAndSale := sales_and_bills.SalesAndBill{
		SaleID:     saleID,
		BillNumber: billNumber,
		Path:       path,
	}
	return db.Create(&billAndSale).Error
}

func GetBillFileName(billNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&sales_and_bills.SalesAndBill{}).
		Where("bill_number = ?", billNumber).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func GetOfferFileName(offerNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&sales_and_offers.SalesAndOffer{}).
		Where("offer_number = ?", offerNumber).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func GetContractFileName(contractNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&sales_and_contracts.SalesAndContract{}).
		Where("contract_number = ?", contractNumber).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func GetDocumentFileName(documentName string) (string, error) {
	var fileName string
	if err := database.DB.Model(&sales_and_documents.SalesAndDocument{}).
		Where("name = ?", documentName).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func DeleteOrderBill(billNumber string) error {
	return database.DB.Model(&sales_and_bills.SalesAndBill{}).
		Where("bill_number = ?", billNumber).
		Delete(&sales_and_bills.SalesAndBill{}).
		Error
}

func DeleteOrderOffer(offerNumber string) error {
	return database.DB.Model(&sales_and_offers.SalesAndOffer{}).
		Where("offer_number = ?", offerNumber).
		Delete(&sales_and_offers.SalesAndOffer{}).
		Error
}

func DeleteOrderContract(contractNumber string) error {
	return database.DB.Model(&sales_and_contracts.SalesAndContract{}).
		Where("contract_number = ?", contractNumber).
		Delete(&sales_and_contracts.SalesAndContract{}).
		Error
}

func DeleteOrderDocument(documentName string) error {
	return database.DB.Model(&sales_and_documents.SalesAndDocument{}).
		Where("name = ?", documentName).
		Delete(&sales_and_documents.SalesAndDocument{}).
		Error
}

func DeleteAllOrderProducts(db *gorm.DB, saleID int64) error {
	return db.
		Where("sale_id = ?", saleID).
		Delete(&sales_and_products.SalesAndProduct{}).Error
}

func GetAllProducts() ([]products.Product, error) {
	var products []products.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func DeleteAllGateOptions(db *gorm.DB, saleID int64, gateID int64) error {
	return db.
		Where("sale_id = ? AND row_number = ?", saleID, gateID).
		Delete(&gates_and_sales_options.GatesAndSalesOption{}).Error
}

func GetIndustrialDriveById(db *gorm.DB, driveID int64) (industrial_gate_drives.IndustrialGateDrive, error) {
	var industrialDrive industrial_gate_drives.IndustrialGateDrive
	if err := db.Where("id = ?", driveID).First(&industrialDrive).Error; err != nil {
		return industrial_gate_drives.IndustrialGateDrive{}, err
	}

	return industrialDrive, nil
}

func GetResidentialDriveById(db *gorm.DB, driveID int64) (residential_gate_drives.ResidentialGateDrive, error) {
	var residentialDrive residential_gate_drives.ResidentialGateDrive
	if err := db.Where("id = ?", driveID).First(&residentialDrive).Error; err != nil {
		return residential_gate_drives.ResidentialGateDrive{}, err
	}

	return residentialDrive, nil
}

func GetRailById(db *gorm.DB, railID int64) (rails.Rail, error) {
	var rail rails.Rail
	if err := db.Where("id = ?", railID).First(&rail).Error; err != nil {
		return rails.Rail{}, err
	}

	return rail, nil
}

func GetOptionsByIDs(tx *gorm.DB, optionIDs []int64) ([]options.Option, error) {
	var optionsList []options.Option

	err := tx.
		Where("id IN ?", optionIDs).
		Find(&optionsList).
		Error

	if err != nil {
		return nil, err
	}

	return optionsList, nil
}

func GetLiftTypeById(db *gorm.DB, liftTypeId int64) (lift_types.LiftType, error) {
	var liftType lift_types.LiftType
	if err := db.Where("id = ?", liftTypeId).First(&liftType).Error; err != nil {
		return lift_types.LiftType{}, err
	}

	return liftType, nil
}

func GetCycleAmountById(db *gorm.DB, cycleAmountId int64) (cycle_amount.CycleAmount, error) {
	var cycleAmount cycle_amount.CycleAmount
	if err := db.Where("id = ?", cycleAmountId).First(&cycleAmount).Error; err != nil {
		return cycle_amount.CycleAmount{}, err
	}

	return cycleAmount, nil
}
