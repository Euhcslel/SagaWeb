// Package database содержит функцию инциализации базы данных и глобальную переменную DB.
package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
)

var DB *gorm.DB

func InitDB() error {
	dsn := os.Getenv("DB_CONNECTION")
	if dsn == "" {
		return errs.ErrEnvVariableNotFound
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Println("База данных подключена успешно")
	return nil
}
