package repository

import (
	"project/internal/database"
	"project/internal/domain/dealers"
)

func GetDealersList() ([]dealers.Dealer, error) {
	var dealersList []dealers.Dealer
	err := database.DB.Preload("User").Preload("Company").Find(&dealersList).Error
	return dealersList, err
}
