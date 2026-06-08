package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_bills"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_appendices"
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

func AttachAppendiceToOrder(db *gorm.DB, orderID int64, appendixNumber string, path string) error {
	return db.Create(&order_appendices.OrderAppendix{
		OrderID:        orderID,
		AppendiceNumber: appendixNumber,
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

func GetAppendiceFileName(orderID int64, appendixNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_appendices.OrderAppendix{}).
		Where("appendix_number = ? AND order_id = ?", appendixNumber, orderID).
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

func DeleteOrderAppendix(db *gorm.DB, appendixNumber string) error {
	return db.
		Where("appendix_number = ?", appendixNumber).
		Delete(&order_appendices.OrderAppendix{}).
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

func GetAppendicesNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var appendicesNumberList []string
	if err := db.Model(&order_appendices.OrderAppendix{}).
		Where("order_id = ?", orderID).
		Pluck("appendix_number", &appendicesNumberList).Error; err != nil {
		return nil, err
	}

	return appendicesNumberList, nil
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
