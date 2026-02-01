package database

import (
	"log"
	"os"
	"project/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var NamedModels map[string]any

func InitDB() *gorm.DB {
	dsn := os.Getenv("DB_CONNECTION")
	if dsn == "" {
		log.Fatal("DB_CONNECTION не найден в .env файле")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	NamedModels = map[string]any{
        "roles": models.Role{},
		"colors": models.Color{},
		"companies": models.Company{},
		"cycle_amount": models.CycleAmount{},
		"dealers": models.Dealer{},
		"gate_type": models.GateType{},
		"gates_and_sales_manual_drive": models.GatesAndSalesManualDrive{},
		"gates_and_sales_option": models.GatesAndSalesOption{},
		"industrial_gate_drives": models.IndustrialGateDrive{},
		"industrial_gates_and_sales_drive": models.IndustrialGatesAndSalesDrive{},
		"lift_types": models.LiftType{},
		"options": models.Option{},
		"products": models.Product{},
		"rails": models.Rail{},
		"residential_gate_drives": models.ResidentialGateDrive{},
		"residential_gates_and_sales_drive_rail": models.ResidentialGatesAndSalesDriveRail{},
		"sales_and_bills": models.SalesAndBill{},
		"sales_and_contracts": models.SalesAndContract{},
		"sales_and_documents": models.SalesAndDocument{},
		"sales_and_gates": models.SalesAndGate{},
		"sales_and_offers": models.SalesAndOffer{},
		"sales_and_products": models.SalesAndProduct{},
		"sales": models.Sale{},
		"sessions": models.Session{},
		"standart_equipment": models.StandartEquipment{},
		"statuses": models.Status{},
		"units": models.Unit{},
		"users": models.User{},
		"sizes": models.Size{},
		"managers_and_dealers": models.ManagerAndDealer{},
		"dealer_reg_requests": models.DealerRegRequest{},
    }

	log.Println("База данных подключена успешно")
	return DB
}

func AutoMigrateAll(db *gorm.DB) error {
	modelsToMigrate := []any{
		&models.Role{},
		&models.Color{},
		&models.Company{},
		&models.CycleAmount{},
		&models.Dealer{},
		&models.GateType{},
		&models.GatesAndSalesManualDrive{},
		&models.GatesAndSalesOption{},
		&models.IndustrialGateDrive{},
		&models.IndustrialGatesAndSalesDrive{},
		&models.LiftType{},
		&models.Option{},
		&models.Product{},
		&models.Rail{},
		&models.ResidentialGateDrive{},
		&models.ResidentialGatesAndSalesDriveRail{},
		&models.SalesAndBill{},
		&models.SalesAndContract{},
		&models.SalesAndDocument{},
		&models.SalesAndGate{},
		&models.SalesAndOffer{},
		&models.SalesAndProduct{},
		&models.Sale{},
		&models.Session{},
		&models.StandartEquipment{},
		&models.Status{},
		&models.Unit{},
		&models.User{},
		&models.Size{},
		&models.ManagerAndDealer{},
		&models.DealerRegRequest{},
	}

	err := db.AutoMigrate(modelsToMigrate...)
	if err != nil {
		return err
	}

	return nil
}
