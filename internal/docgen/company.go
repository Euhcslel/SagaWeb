package docgen

import (
	"os"
	"strconv"
)

// CompanyInfo хранит реквизиты компании-поставщика, используемые во всех документах.
type CompanyInfo struct {
	Name               string // Название, напр. «ООО «САГА»»
	INN                string // ИНН
	KPP                string // КПП
	WarrantyEquipment  string // Гарантия на оборудование, напр. «12 месяцев»
	WarrantyAutomation string // Гарантия на автоматику, напр. «24 месяца»
	OfferValidityDays  int    // Срок действия КП в днях
}

// SupplierInfo содержит данные поставщика (дилера) для шапки КП.
type SupplierInfo struct {
	Name  string // Название компании дилера
	Phone string // Телефон
	Email string // Email
}

var defaultCompany CompanyInfo

// InitCompany читает реквизиты из переменных окружения и сохраняет их.
// Вызывать один раз при старте приложения после godotenv.Load().
func InitCompany() {
	days, _ := strconv.Atoi(os.Getenv("OFFER_VALIDITY_DAYS"))
	if days <= 0 {
		days = 14
	}
	defaultCompany = CompanyInfo{
		Name:               os.Getenv("COMPANY_NAME"),
		INN:                os.Getenv("COMPANY_INN"),
		KPP:                os.Getenv("COMPANY_KPP"),
		WarrantyEquipment:  os.Getenv("WARRANTY_EQUIPMENT"),
		WarrantyAutomation: os.Getenv("WARRANTY_AUTOMATION"),
		OfferValidityDays:  days,
	}
}

// CompanyInfoFromEnv возвращает реквизиты, инициализированные при старте.
func CompanyInfoFromEnv() CompanyInfo {
	return defaultCompany
}
