package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_bills"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_contracts"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_documents"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_offers"
	"gorm.io/gorm"
)

func AttachOfferToOrder(db *gorm.DB, orderID int64, offerNumber string, path string) error {
	return db.Create(&order_offers.OrderOffer{
		OrderID:     orderID,
		OfferNumber: offerNumber,
		Path:        path,
	}).Error
}

func AttachContractToOrder(db *gorm.DB, orderID int64, contractNumber string, path string) error {
	return db.Create(&order_contracts.OrderContract{
		OrderID:        orderID,
		ContractNumber: contractNumber,
		Path:           path,
	}).Error
}

func AttachBillToOrder(db *gorm.DB, orderID int64, billNumber string, path string) error {
	return db.Create(&order_bills.OrderBill{
		OrderID:    orderID,
		BillNumber: billNumber,
		Path:       path,
	}).Error
}

func GetBillFileName(orderID int64, billNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_bills.OrderBill{}).
		Where("bill_number = ? AND order_id = ?", billNumber, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func GetOfferFileName(orderID int64, offerNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_offers.OrderOffer{}).
		Where("offer_number = ? AND order_id = ?", offerNumber, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func GetContractFileName(orderID int64, contractNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_contracts.OrderContract{}).
		Where("contract_number = ? AND order_id = ?", contractNumber, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func GetDocumentFileName(orderID int64, documentName string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_documents.OrderDocument{}).
		Where("name = ? AND order_id = ?", documentName, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

func DeleteOrderBill(db *gorm.DB, billNumber string) error {
	return db.
		Where("bill_number = ?", billNumber).
		Delete(&order_bills.OrderBill{}).
		Error
}

func DeleteOrderOffer(db *gorm.DB, offerNumber string) error {
	return db.
		Where("offer_number = ?", offerNumber).
		Delete(&order_offers.OrderOffer{}).
		Error
}

func DeleteOrderContract(db *gorm.DB, contractNumber string) error {
	return db.
		Where("contract_number = ?", contractNumber).
		Delete(&order_contracts.OrderContract{}).
		Error
}

func DeleteOrderDocument(db *gorm.DB, documentName string) error {
	return db.
		Where("name = ?", documentName).
		Delete(&order_documents.OrderDocument{}).
		Error
}

func GetDocumentsNameList(db *gorm.DB, orderID int64) ([]string, error) {
	var documentsNameList []string
	if err := db.Model(&order_documents.OrderDocument{}).
		Where("order_id = ?", orderID).
		Pluck("name", &documentsNameList).Error; err != nil {
		return nil, err
	}

	return documentsNameList, nil
}

func GetOffersNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var offersNumberList []string
	if err := db.Model(&order_offers.OrderOffer{}).
		Where("order_id = ?", orderID).
		Pluck("offer_number", &offersNumberList).Error; err != nil {
		return nil, err
	}

	return offersNumberList, nil
}

func GetContractsNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var contractsNumberList []string
	if err := db.Model(&order_contracts.OrderContract{}).
		Where("order_id = ?", orderID).
		Pluck("contract_number", &contractsNumberList).Error; err != nil {
		return nil, err
	}

	return contractsNumberList, nil
}

func GetBillsNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var billsNumberList []string
	if err := db.Model(&order_bills.OrderBill{}).
		Where("order_id = ?", orderID).
		Pluck("bill_number", &billsNumberList).Error; err != nil {
		return nil, err
	}

	return billsNumberList, nil
}

func GetManualDrivePrice(db *gorm.DB) (*manual_drive_prices.ManualDrivePrice, error) {
	var manualDrivePrice manual_drive_prices.ManualDrivePrice
	if err := db.First(&manualDrivePrice).Error; err != nil {
		return nil, err
	}
	return &manualDrivePrice, nil
}
