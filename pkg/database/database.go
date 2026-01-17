package database

import (
	"log"
	"os"
	"project/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
		&models.MontageType{},
		&models.MontagesAndSale{},
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
