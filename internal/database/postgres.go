package database

import (
	"log"
	"os"

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
