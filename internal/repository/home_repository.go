package repository

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealers"
)

func GetDealersList() ([]dealers.Dealer, error) {
	var dealersList []dealers.Dealer
	if err := database.DB.Preload("User").Preload("Company").Find(&dealersList).Error; err != nil {
		return nil, err
	}
	
	return dealersList, nil
}
