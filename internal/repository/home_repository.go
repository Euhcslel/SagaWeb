package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
)

func GetDealersList() ([]dealers.Dealer, error) {
	var dealersList []dealers.Dealer
	err := database.DB.Preload("User").Preload("Company").Find(&dealersList).Error
	return dealersList, err
}
