// Package docgen генерирует документы для заказов: коммерческие предложения (PDF),
// приложения к договорам (PDF) и договоры (DOCX) на основе шаблонов и данных заказа.

package docgen

import (
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"os"
	"strconv"
)

// SupplierInfo содержит данные поставщика (дилера) для шапки КП.
type SupplierInfo struct {
	Name  string
	Phone string
	Email string
}

// CustomerInfo содержит данные заказчика для приложения к договору.
type CustomerInfo struct {
	OrganizationName string
	ContactPerson    string
	Phone            string
}

// CompanyConfig содержит стандартные реквизиты компании, а также некоторые необходимые параметры.
// OfferValidityDays - срок действия коммерческого предложения.
// WarrantyEquip - гарантия на оборудование.
// WarrantyAuto - гарантия на автоматику.
// CompanyName - наименование компании.
// CompanyINN - ИНН компании.
// CompanyKPP - КПП компании.
// Representative - ФИО руководителя.
type CompanyConfig struct {
	OfferValidityDays    int
	WarrantyEquip        string
	WarrantyAuto         string
	CompanyName          string
	CompanyINN           string
	CompanyKPP           string
	Representative       string
	ContractTemplatePath string
}

var Cfg CompanyConfig

// Init читает настройки из переменных окружения.
func Init() error {
	days, err := strconv.Atoi(os.Getenv("OFFER_VALIDITY_DAYS"))
	if days <= 0 {
		days = 14
	}
	if err != nil {
		return errs.ErrParseStringToInt
	}
	Cfg.OfferValidityDays = days
	Cfg.WarrantyEquip = os.Getenv("WARRANTY_EQUIPMENT")
	Cfg.WarrantyAuto = os.Getenv("WARRANTY_AUTOMATION")
	Cfg.CompanyName = os.Getenv("COMPANY_NAME")
	Cfg.CompanyINN = os.Getenv("COMPANY_INN")
	Cfg.CompanyKPP = os.Getenv("COMPANY_KPP")
	Cfg.Representative = os.Getenv("COMPANY_REPRESENTATIVE")
	Cfg.ContractTemplatePath = os.Getenv("CONTRACT_TEMPLATE_PATH")

	return nil
}
