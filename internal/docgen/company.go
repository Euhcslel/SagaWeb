package docgen

import (
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

var (
	defaultOfferValidity   int
	defaultWarrantyEquip   string
	defaultWarrantyAuto    string
	defaultCompanyName     string
	defaultCompanyINN      string
	defaultCompanyKPP      string
	defaultRepresentative  string
)

// Init читает настройки из переменных окружения.
// Вызывать один раз при старте приложения после godotenv.Load().
func Init() {
	days, _ := strconv.Atoi(os.Getenv("OFFER_VALIDITY_DAYS"))
	if days <= 0 {
		days = 14
	}
	defaultOfferValidity  = days
	defaultWarrantyEquip  = os.Getenv("WARRANTY_EQUIPMENT")
	defaultWarrantyAuto   = os.Getenv("WARRANTY_AUTOMATION")
	defaultCompanyName    = os.Getenv("COMPANY_NAME")
	defaultCompanyINN     = os.Getenv("COMPANY_INN")
	defaultCompanyKPP     = os.Getenv("COMPANY_KPP")
	defaultRepresentative = os.Getenv("COMPANY_REPRESENTATIVE")
}

// offerValidityDays возвращает срок действия КП по умолчанию.
func offerValidityDays() int { return defaultOfferValidity }
