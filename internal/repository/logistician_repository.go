package repository

import (
	"time"

	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/gate_sizes"
	"github.com/Euhcslel/SagaWeb/internal/domain/orders"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func GetAllOrders() ([]orders.Order, error) {
	var result []orders.Order
	if err := database.DB.
		Preload("Dealer.User").
		Preload("Manager").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateOrderManufactureDate(orderID int64, manufactureDate *time.Time) error {
	return database.DB.Model(&orders.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]any{"manufacture_date": manufactureDate}).Error
}

func UpdateGateSizePrice(width, height int64, gateType enums.GateType, wholesalePrice, retailPrice decimal.Decimal) error {
	size := gate_sizes.GateSize{
		Width:          width,
		Height:         height,
		GateType:       gateType,
		WholesalePrice: wholesalePrice,
		RetailPrice:    retailPrice,
	}
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "width"}, {Name: "height"}, {Name: "gate_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"wholesale_price", "retail_price"}),
	}).Create(&size).Error
}
