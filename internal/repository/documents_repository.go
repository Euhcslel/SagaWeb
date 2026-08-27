package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_bills"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_appendices"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_documents"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_offers"
	"gorm.io/gorm"
)

// AttachOfferToOrder прикрепляет коммерческое предложение к заказу в базе данных.
func AttachOfferToOrder(db *gorm.DB, orderID int64, offerNumber string, path string) error {
	return db.Create(&order_offers.OrderOffer{
		OrderID:     orderID,
		OfferNumber: offerNumber,
		Path:        path,
	}).Error
}

// AttachAppendiceToOrder прикрепляет приложение к договору к заказу в базе данных.
func AttachAppendiceToOrder(db *gorm.DB, orderID int64, appendiceNumber string, path string) error {
	return db.Create(&order_appendices.OrderAppendix{
		OrderID:        orderID,
		AppendiceNumber: appendiceNumber,
		Path:           path,
	}).Error
}

// AttachBillToOrder прикрепляет счет к заказу в базе данных.
func AttachBillToOrder(db *gorm.DB, orderID int64, billNumber string, path string) error {
	return db.Create(&order_bills.OrderBill{
		OrderID:    orderID,
		BillNumber: billNumber,
		Path:       path,
	}).Error
}

// AttachDocumentToOrder прикрепляет документ к заказу в базе данных.
func AttachDocumentToOrder(db *gorm.DB, orderID int64, name string, path string) error {
	return db.Create(&order_documents.OrderDocument{
		OrderID: orderID,
		Name:    name,
		Path:    path,
	}).Error
}

// GetBillFileName возвращает название файла счета по номеру счета и идентификатору заказа.
func GetBillFileName(orderID int64, billNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_bills.OrderBill{}).
		Where("bill_number = ? AND order_id = ?", billNumber, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

// GetOfferFileName возвращает название файла коммерческого предложения по номеру предложения и идентификатору заказа.
func GetOfferFileName(orderID int64, offerNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_offers.OrderOffer{}).
		Where("offer_number = ? AND order_id = ?", offerNumber, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

// GetAppendiceFileName возвращает название файла приложения к договору по номеру приложения и идентификатору заказа.
func GetAppendiceFileName(orderID int64, appendiceNumber string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_appendices.OrderAppendix{}).
		Where("appendice_number = ? AND order_id = ?", appendiceNumber, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

// GetDocumentFileName возвращает название файла документа по имени документа и идентификатору заказа.
func GetDocumentFileName(orderID int64, documentName string) (string, error) {
	var fileName string
	if err := database.DB.Model(&order_documents.OrderDocument{}).
		Where("name = ? AND order_id = ?", documentName, orderID).
		Pluck("path", &fileName).Error; err != nil {
		return "", err
	}

	return fileName, nil
}

// DeleteOrderBill удаляет запись о счете из базы данных по номеру счета и идентификатору заказа.
func DeleteOrderBill(db *gorm.DB, orderID int64, billNumber string) error {
	return db.
		Where("bill_number = ? AND order_id = ?", billNumber, orderID).
		Delete(&order_bills.OrderBill{}).
		Error
}

// DeleteOrderOffer удаляет запись о коммерческом предложении из базы данных по номеру предложения и идентификатору заказа.
func DeleteOrderOffer(db *gorm.DB, orderID int64, offerNumber string) error {
	return db.
		Where("offer_number = ? AND order_id = ?", offerNumber, orderID).
		Delete(&order_offers.OrderOffer{}).
		Error
}

// DeleteOrderAppendice удаляет запись о приложении к договору из базы данных по номеру приложения и идентификатору заказа.
func DeleteOrderAppendice(db *gorm.DB, orderID int64, appendiceNumber string) error {
	return db.
		Where("appendice_number = ? AND order_id = ?", appendiceNumber, orderID).
		Delete(&order_appendices.OrderAppendix{}).
		Error
}

// DeleteOrderDocument удаляет запись о документе из базы данных по имени документа и идентификатору заказа.
func DeleteOrderDocument(db *gorm.DB, orderID int64, documentName string) error {
	return db.
		Where("name = ? AND order_id = ?", documentName, orderID).
		Delete(&order_documents.OrderDocument{}).
		Error
}

// GetDocumentsNameList возвращает список имен документов, связанных с указанным идентификатором заказа.
func GetDocumentsNameList(db *gorm.DB, orderID int64) ([]string, error) {
	var documentsNameList []string
	if err := db.Model(&order_documents.OrderDocument{}).
		Where("order_id = ?", orderID).
		Pluck("name", &documentsNameList).Error; err != nil {
		return nil, err
	}

	return documentsNameList, nil
}

// GetOffersNumberList возвращает список номеров коммерческих предложений, связанных с указанным идентификатором заказа.
func GetOffersNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var offersNumberList []string
	if err := db.Model(&order_offers.OrderOffer{}).
		Where("order_id = ?", orderID).
		Pluck("offer_number", &offersNumberList).Error; err != nil {
		return nil, err
	}

	return offersNumberList, nil
}

// GetAppendicesNumberList возвращает список номеров приложений к договорам, связанных с указанным идентификатором заказа.
func GetAppendicesNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var appendicesNumberList []string
	if err := db.Model(&order_appendices.OrderAppendix{}).
		Where("order_id = ?", orderID).
		Pluck("appendice_number", &appendicesNumberList).Error; err != nil {
		return nil, err
	}

	return appendicesNumberList, nil
}

// GetBillsNumberList возвращает список номеров счетов, связанных с указанным идентификатором заказа.
func GetBillsNumberList(db *gorm.DB, orderID int64) ([]string, error) {
	var billsNumberList []string
	if err := db.Model(&order_bills.OrderBill{}).
		Where("order_id = ?", orderID).
		Pluck("bill_number", &billsNumberList).Error; err != nil {
		return nil, err
	}

	return billsNumberList, nil
}
