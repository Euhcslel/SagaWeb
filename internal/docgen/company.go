package docgen

import (
	"os"
	"strconv"
)

// CompanyInfo хранит реквизиты компании-производителя для приложений к договорам.
type CompanyInfo struct {
	Name               string
	INN                string
	KPP                string
	WarrantyEquipment  string
	WarrantyAutomation string
}

// SupplierInfo содержит данные поставщика (дилера) для шапки КП.
type SupplierInfo struct {
	Name              string
	Phone             string
	Email             string
	OfferValidityDays int
}

var (
	defaultCompany        CompanyInfo
	defaultOfferValidity  int
)

// InitCompany читает реквизиты из переменных окружения и сохраняет их.
// Вызывать один раз при старте приложения после godotenv.Load().
func InitCompany() {
	days, _ := strconv.Atoi(os.Getenv("OFFER_VALIDITY_DAYS"))
	if days <= 0 {
		days = 14
	}
	defaultOfferValidity = days
	defaultCompany = CompanyInfo{
		Name:               os.Getenv("COMPANY_NAME"),
		INN:                os.Getenv("COMPANY_INN"),
		KPP:                os.Getenv("COMPANY_KPP"),
		WarrantyEquipment:  os.Getenv("WARRANTY_EQUIPMENT"),
		WarrantyAutomation: os.Getenv("WARRANTY_AUTOMATION"),
	}
}

// CompanyInfoFromEnv возвращает реквизиты производителя, инициализированные при старте.
func CompanyInfoFromEnv() CompanyInfo {
	return defaultCompany
}

// OfferValidityDaysFromEnv возвращает срок действия КП по умолчанию.
func OfferValidityDaysFromEnv() int {
	return defaultOfferValidity
}
