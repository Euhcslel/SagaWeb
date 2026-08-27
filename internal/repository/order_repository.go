package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/dealer_manager_assignments"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gate_industrial_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gate_manual_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gate_options"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gate_residential_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_products"
	"github.com/Euhcslel/SagaWeb/internal/domain/orders"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"

	"github.com/Euhcslel/SagaWeb/internal/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetOrderWithCustomer возвращает заказ с информацией о клиенте.
func GetOrderWithCustomer(db *gorm.DB, orderID int64) (*orders.Order, error) {
	var order orders.Order
	if err := db.
		Preload("Dealer.User").
		Preload("Dealer.Company").
		Preload("Manager").
		Where("id = ?", orderID).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAllDealerOrders возвращает все заказы, связанные с указанным дилером.
func GetAllDealerOrders(db *gorm.DB, user *users.User) ([]orders.Order, error) {
	var orders []orders.Order
	if err := db.
		Preload("Dealer.User").
		Preload("Manager").
		Where("dealer_id = ?", user.ID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetAllManagerOrders возвращает все заказы, связанные с указанным менеджером.
func GetAllManagerOrders(db *gorm.DB, user *users.User) ([]orders.Order, error) {
	var orders []orders.Order
	if err := db.
		Preload("Dealer.User").
		Preload("Manager").
		Where("manager_id = ?", user.ID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderGatesByOrderID возвращает все ворота, связанные с указанным заказом.
func GetOrderGatesByOrderID(db *gorm.DB, orderID int64) ([]order_gates.OrderGate, error) {
	var orderGates []order_gates.OrderGate
	if err := db.
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("order_id = ?", orderID).
		Find(&orderGates).Error; err != nil {
		return nil, err
	}

	return orderGates, nil
}

// GetAllOptionsForOrder возвращает все дополнительные опции, которые были выбраны для ворот в указанном заказе.
func GetAllOptionsForOrder(db *gorm.DB, orderID int64, gateIDs []int64) ([]order_gate_options.OrderGateOption, error) {
	var gateOptions []order_gate_options.OrderGateOption
	if err := db.
		Where("order_id = ? AND row_number IN ?", orderID, gateIDs).
		Preload("Option").
		Find(&gateOptions).Error; err != nil {
		return nil, err
	}

	return gateOptions, nil
}

// GetAllOrderProducts возвращает все товары, связанные с указанным заказом.
func GetAllOrderProducts(db *gorm.DB, orderID int64) ([]order_products.OrderProduct, error) {
	var products []order_products.OrderProduct
	if err := db.Preload("Product").Where("order_id = ?", orderID).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// GetCurrentGate возвращает информацию о конкретных воротах в указанном заказе.
func GetCurrentGate(db *gorm.DB, orderID int64, gateID int64) (*order_gates.OrderGate, error) {
	var gate order_gates.OrderGate
	if err := db.Model(order_gates.OrderGate{}).
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("order_id = ? and row_number = ?", orderID, gateID).
		First(&gate).Error; err != nil {
		return nil, err
	}

	return &gate, nil
}

// GetCurrentGateOptions возвращает список всех дополнительных опций,
// связанных с конкретными воротами в указанном заказе.
func GetCurrentGateOptions(db *gorm.DB, orderID int64, gateID int64) ([]order_gate_options.OrderGateOption, error) {
	var options []order_gate_options.OrderGateOption
	if err := db.
		Model(order_gate_options.OrderGateOption{}).
		Preload("Option").Where("order_id = ? and row_number = ?", orderID, gateID).
		Find(&options).Error; err != nil {
		return nil, err
	}

	return options, nil
}

// GetIndustrialDriveForGate возвращает информацию о промышленном приводе для конкретных ворот в указанном заказе.
func GetIndustrialDriveForGate(db *gorm.DB, orderID int64, gateID int64) (*order_gate_industrial_drives.OrderGateIndustrialDrive, error) {
	var drive order_gate_industrial_drives.OrderGateIndustrialDrive
	if err := db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		First(&drive).Error; err != nil {
		return nil, err
	}

	return &drive, nil
}

// GetResidentialDriveForGate возвращает информацию о бытовом приводе для конкретных ворот в указанном заказе.
func GetResidentialDriveForGate(db *gorm.DB, orderID int64, gateID int64) (*order_gate_residential_drives.OrderGateResidentialDrive, error) {
	var drive order_gate_residential_drives.OrderGateResidentialDrive
	if err := db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		First(&drive).Error; err != nil {
		return nil, err
	}

	return &drive, nil
}

// GetManualDriveForGate возвращает информацию о ручном приводе для конкретных ворот в указанном заказе.
func GetManualDriveForGate(db *gorm.DB, orderID int64, gateID int64) (*order_gate_manual_drives.OrderGateManualDrive, error) {
	var drive order_gate_manual_drives.OrderGateManualDrive
	if err := db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		First(&drive).Error; err != nil {
		return nil, err
	}

	return &drive, nil
}

// DeleteOrder удаляет заказ из базы данных по его идентификатору.
func DeleteOrder(db *gorm.DB, orderID int64) error {
	return db.Where("id = ?", orderID).Delete(&orders.Order{}).Error
}

// CreateNewGate создает новые ворота.
func CreateNewGate(db *gorm.DB, gate order_gates.OrderGate) (order_gates.OrderGate, error) {
	if err := db.Clauses(clause.Returning{}).Create(&gate).Error; err != nil {
		return order_gates.OrderGate{}, err
	}

	return gate, nil
}

// CreateIndustrialDriveForGate создает запись о промышленном приводе для конкретных ворот в указанном заказе.
func CreateIndustrialDriveForGate(db *gorm.DB, orderID int64, rowNumber int64, driveID int64) error {
	if err := db.Create(&order_gate_industrial_drives.OrderGateIndustrialDrive{
		OrderID:   orderID,
		RowNumber: rowNumber,
		DriveID:   driveID,
	}).Error; err != nil {
		return err
	}

	return nil
}

// CreateResidentialDriveForGate создает запись о бытовом приводе для конкретных ворот в указанном заказе.
func CreateResidentialDriveForGate(db *gorm.DB, orderID int64, rowNumber int64, driveID int64, railID int64) error {
	if err := db.Create(&order_gate_residential_drives.OrderGateResidentialDrive{
		OrderID:   orderID,
		RowNumber: rowNumber,
		DriveID:   driveID,
		RailID:    railID,
	}).Error; err != nil {
		return err
	}

	return nil
}

// CreateManualDriveForGate создает запись о ручном приводе для конкретных ворот в указанном заказе.
func CreateManualDriveForGate(db *gorm.DB, orderID int64, rowNumber int64, chainLength int32) error {
	if err := db.Create(&order_gate_manual_drives.OrderGateManualDrive{
		OrderID:     orderID,
		RowNumber:   rowNumber,
		ChainLength: chainLength,
	}).Error; err != nil {
		return err
	}

	return nil
}

// GetManagerIDByDealerID возвращает идентификатор менеджера, связанного с указанным дилером.
func GetManagerIDByDealerID(db *gorm.DB, dealerID int64) (int64, error) {
	var managerID int64
	if err := db.Model(&dealer_manager_assignments.DealerManagerAssignment{}).
		Where("dealer_id = ?", dealerID).
		Pluck("manager_id", &managerID).Error; err != nil {
		log.Println(err)

		return 0, err
	}

	return managerID, nil
}

// CreateNewOrder создает новый заказ в базе данных.
func CreateNewOrder(db *gorm.DB, order *orders.Order) error {
	return db.Create(&order).Error
}

// CreateOrderProducts создает записи о товарах для указанного заказа.
func CreateOrderProducts(db *gorm.DB, products []order_products.OrderProduct) error {
	return db.Create(&products).Error
}

// CreateGateOptions создает записи о дополнительных опциях для конкретных ворот в заказе.
func CreateGateOptions(db *gorm.DB, options []order_gate_options.OrderGateOption) error {
	return db.Create(&options).Error
}

// DeleteGateFromOrder удаляет ворота из указанного заказа.
func DeleteGateFromOrder(db *gorm.DB, orderID int64, rowNumber int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ? AND row_number = ?", orderID, rowNumber).
			Delete(&order_gates.OrderGate{}).Error; err != nil {
			return err
		}

		var remaining []order_gates.OrderGate
		if err := tx.Where("order_id = ? AND row_number > ?", orderID, rowNumber).
			Order("row_number ASC").
			Find(&remaining).Error; err != nil {
			return err
		}

		for _, g := range remaining {
			if err := tx.Model(&order_gates.OrderGate{}).
				Where("order_id = ? AND row_number = ?", orderID, g.RowNumber).
				Update("row_number", g.RowNumber-1).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteGateResidentialDrive удаляет запись о бытовом приводе для конкретных ворот в заказе.
func DeleteGateResidentialDrive(db *gorm.DB, orderID int64, gateID int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(&order_gate_residential_drives.OrderGateResidentialDrive{}).
		Error
}

// DeleteGateIndustrialDrive удаляет запись о промышленном приводе для конкретных ворот в заказе.
func DeleteGateIndustrialDrive(db *gorm.DB, orderID int64, gateID int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(order_gate_industrial_drives.OrderGateIndustrialDrive{}).
		Error
}

// DeleteGateManualDrive удаляет запись о ручном приводе для конкретных ворот в заказе.
func DeleteGateManualDrive(db *gorm.DB, orderID int64, gateID int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(&order_gate_manual_drives.OrderGateManualDrive{}).Error
}

// UpdateGateIndustrialDrive обновляет запись о промышленном приводе для конкретных ворот в заказе.
func UpdateGateIndustrialDrive(db *gorm.DB, orderID int64, gateID int64, driveID int64) error {
	return db.Model(&order_gate_industrial_drives.OrderGateIndustrialDrive{}).
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Update("drive_id", driveID).Error
}

// UpdateGateResidentialDrive обновляет запись о бытовом приводе для конкретных ворот в заказе.
func UpdateGateResidentialDrive(db *gorm.DB, orderID int64, gateID int64, driveID int64, railID int64) error {
	return db.Model(&order_gate_residential_drives.OrderGateResidentialDrive{}).
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Updates(map[string]any{"drive_id": driveID, "rail_id": railID}).Error
}

// UpdateGateManualDrive обновляет запись о ручном приводе для конкретных ворот в заказе.
func UpdateGateManualDrive(db *gorm.DB, orderID int64, gateID int64, chainLength int32) error {
	return db.Model(&order_gate_manual_drives.OrderGateManualDrive{}).
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Update("chain_length", chainLength).Error
}

// UpdateGate обновляет запись о воротах в заказе.
func UpdateGate(db *gorm.DB, gate order_gates.OrderGate) error {
	return db.Save(&gate).Error
}

// UpdateOrderStatus обновляет статус заказа в базе данных.
func UpdateOrderStatus(db *gorm.DB, orderID int64, status enums.OrderStatus) error {
	return db.Model(&orders.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// SetOrderFinalizedAt устанавливает время завершения заказа в базе данных.
func SetOrderFinalizedAt(db *gorm.DB, orderID int64) error {
	return db.Model(&orders.Order{}).
		Where("id = ? AND finalized_at IS NULL", orderID).
		Update("finalized_at", db.Raw("CURRENT_TIMESTAMP")).Error
}

// GetOrderFinalizedAt возвращает время завершения заказа из базы данных.
func GetOrderFinalizedAt(db *gorm.DB, orderID int64) (*time.Time, error) {
	var t *time.Time
	err := db.Model(&orders.Order{}).Where("id = ?", orderID).Pluck("finalized_at", &t).Error
	return t, err
}

// GetOrderManufactureDate возвращает дату производства заказа из базы данных.
func GetOrderManufactureDate(db *gorm.DB, orderID int64) (*time.Time, error) {
	var nt sql.NullTime
	err := db.Model(&orders.Order{}).Where("id = ?", orderID).Pluck("manufacture_date", &nt).Error
	if err != nil {
		return nil, err
	}
	if !nt.Valid {
		return nil, nil
	}
	return &nt.Time, nil
}

// GetOrderStatus возвращает статус заказа из базы данных.
func GetOrderStatus(orderID int64) (*string, error) {
	var status string
	if err := database.DB.Model(&orders.Order{}).
		Where("id = ?", orderID).
		Pluck("status", &status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

// DeleteAllOrderProducts удаляет все товары, связанные с указанным заказом.
func DeleteAllOrderProducts(db *gorm.DB, orderID int64) error {
	return db.
		Where("order_id = ?", orderID).
		Delete(&order_products.OrderProduct{}).Error
}

// GetAllProducts возвращает все товары из базы данных.
func GetAllProducts() ([]products.Product, error) {
	var products []products.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// DeleteAllGateOptions удаляет все дополнительные опции, связанные с конкретными воротами в заказе.
func DeleteAllGateOptions(db *gorm.DB, orderID int64, gateID int64) error {
	return db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(&order_gate_options.OrderGateOption{}).Error
}

// GetIndustrialDriveByID возвращает информацию о промышленном приводе по его идентификатору.
func GetIndustrialDriveByID(db *gorm.DB, driveID int64) (*industrial_gate_drives.IndustrialGateDrive, error) {
	var industrialDrive industrial_gate_drives.IndustrialGateDrive
	if err := db.Where("id = ?", driveID).First(&industrialDrive).Error; err != nil {
		return nil, err
	}

	return &industrialDrive, nil
}

// GetResidentialDriveByID возвращает информацию о бытовом приводе по его идентификатору.
func GetResidentialDriveByID(db *gorm.DB, driveID int64) (*residential_gate_drives.ResidentialGateDrive, error) {
	var residentialDrive residential_gate_drives.ResidentialGateDrive
	if err := db.Where("id = ?", driveID).First(&residentialDrive).Error; err != nil {
		return nil, err
	}

	return &residentialDrive, nil
}

// GetRailByID возвращает информацию о направляющей по её идентификатору.
func GetRailByID(db *gorm.DB, railID int64) (*rails.Rail, error) {
	var rail rails.Rail
	if err := db.Where("id = ?", railID).First(&rail).Error; err != nil {
		return nil, err
	}

	return &rail, nil
}

// GetOptionsByIDs возвращает список опций по их идентификаторам.
func GetOptionsByIDs(db *gorm.DB, optionIDs []int64) ([]options.Option, error) {
	var optionsList []options.Option
	if err := db.
		Where("id IN ?", optionIDs).
		Find(&optionsList).
		Error; err != nil {
		return nil, err
	}

	return optionsList, nil
}

// GetLiftTypeByID возвращает информацию о типе подъема по его идентификатору.
func GetLiftTypeByID(db *gorm.DB, liftTypeID int64) (*lift_types.LiftType, error) {
	var liftType lift_types.LiftType
	if err := db.Where("id = ?", liftTypeID).First(&liftType).Error; err != nil {
		return nil, err
	}

	return &liftType, nil
}

// GetCycleAmountByID возвращает информацию о количестве циклов по его идентификатору.
func GetCycleAmountByID(db *gorm.DB, cycleAmountID int64) (*cycle_amounts.CycleAmount, error) {
	var cycleAmount cycle_amounts.CycleAmount
	if err := db.Where("id = ?", cycleAmountID).First(&cycleAmount).Error; err != nil {
		return nil, err
	}

	return &cycleAmount, nil
}

// GetColorByID возвращает информацию о цвете по его идентификатору.
func GetColorByID(db *gorm.DB, colorID int64) (*colors.Color, error) {
	var color colors.Color
	if err := db.Where("id = ?", colorID).First(&color).Error; err != nil {
		return nil, err
	}
	return &color, nil
}

// GetProductByID возвращает информацию о товаре по его идентификатору.
func GetProductByID(db *gorm.DB, productID int64) (*products.Product, error) {
	var product products.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductHistoricalPrices вовзращает цены на товары, актуальные на указанную дату.
func GetProductHistoricalPrices(db *gorm.DB, productIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(productIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		ProductID      int64           `gorm:"column:product_id"`
		RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
		WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (product_id) product_id, retail_price, wholesale_price
		FROM product_price_history
		WHERE product_id IN ? AND set_at <= ?
		ORDER BY product_id, set_at DESC
	`, productIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.ProductID] = types.PricePair{RetailPrice: r.RetailPrice, WholesalePrice: r.WholesalePrice}
	}
	return result, err
}

// GetOptionHistoricalPrices возвращает цены на дополнительные опции, актуальные на указанную дату.
func GetOptionHistoricalPrices(db *gorm.DB, optionIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(optionIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		OptionID       int64           `gorm:"column:option_id"`
		RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
		WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (option_id) option_id, retail_price, wholesale_price
		FROM option_price_history
		WHERE option_id IN ? AND set_at <= ?
		ORDER BY option_id, set_at DESC
	`, optionIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.OptionID] = types.PricePair{RetailPrice: r.RetailPrice, WholesalePrice: r.WholesalePrice}
	}
	return result, err
}

// GetIndustrialDriveHistoricalPrices возвращает цены на промышленные приводы, актуальные на указанную дату.
func GetIndustrialDriveHistoricalPrices(db *gorm.DB, driveIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(driveIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		DriveID        int64           `gorm:"column:drive_id"`
		RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
		WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (drive_id) drive_id, retail_price, wholesale_price
		FROM industrial_gate_drive_price_history
		WHERE drive_id IN ? AND set_at <= ?
		ORDER BY drive_id, set_at DESC
	`, driveIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.DriveID] = types.PricePair{RetailPrice: r.RetailPrice, WholesalePrice: r.WholesalePrice}
	}
	return result, err
}

// GetResidentialDriveHistoricalPrices возвращает цены на бытовые приводы, актуальные на указанную дату.
func GetResidentialDriveHistoricalPrices(db *gorm.DB, driveIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(driveIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		DriveID        int64           `gorm:"column:drive_id"`
		RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
		WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (drive_id) drive_id, retail_price, wholesale_price
		FROM residential_gate_drive_price_history
		WHERE drive_id IN ? AND set_at <= ?
		ORDER BY drive_id, set_at DESC
	`, driveIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.DriveID] = types.PricePair{RetailPrice: r.RetailPrice, WholesalePrice: r.WholesalePrice}
	}
	return result, err
}

// GetRailHistoricalPrices возвращает цены на направляющие, актуальные на указанную дату.
func GetRailHistoricalPrices(db *gorm.DB, railIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(railIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		RailID         int64           `gorm:"column:rail_id"`
		RetailPrice    decimal.Decimal `gorm:"column:retail_price"`
		WholesalePrice decimal.Decimal `gorm:"column:wholesale_price"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (rail_id) rail_id, retail_price, wholesale_price
		FROM rail_price_history
		WHERE rail_id IN ? AND set_at <= ?
		ORDER BY rail_id, set_at DESC
	`, railIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.RailID] = types.PricePair{RetailPrice: r.RetailPrice, WholesalePrice: r.WholesalePrice}
	}
	return result, err
}

// GetLiftTypeHistoricalMarkups возвращает наценки на типы подъема, актуальные на указанную дату.
func GetLiftTypeHistoricalMarkups(db *gorm.DB, liftTypeIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(liftTypeIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		LiftTypeID      int64           `gorm:"column:lift_type_id"`
		RetailMarkup    decimal.Decimal `gorm:"column:retail_markup"`
		WholesaleMarkup decimal.Decimal `gorm:"column:wholesale_markup"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (lift_type_id) lift_type_id, retail_markup, wholesale_markup
		FROM lift_type_markup_history
		WHERE lift_type_id IN ? AND set_at <= ?
		ORDER BY lift_type_id, set_at DESC
	`, liftTypeIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.LiftTypeID] = types.PricePair{RetailPrice: r.RetailMarkup, WholesalePrice: r.WholesaleMarkup}
	}
	return result, err
}

// GetCycleAmountHistoricalMarkups возвращает наценки на количество циклов, актуальные на указанную дату.
func GetCycleAmountHistoricalMarkups(db *gorm.DB, cycleAmountIDs []int64, at time.Time) (map[int64]types.PricePair, error) {
	if len(cycleAmountIDs) == 0 {
		return map[int64]types.PricePair{}, nil
	}
	type row struct {
		CycleAmountID   int64           `gorm:"column:cycle_amount_id"`
		RetailMarkup    decimal.Decimal `gorm:"column:retail_markup"`
		WholesaleMarkup decimal.Decimal `gorm:"column:wholesale_markup"`
	}
	var rows []row
	err := db.Raw(`
		SELECT DISTINCT ON (cycle_amount_id) cycle_amount_id, retail_markup, wholesale_markup
		FROM cycle_amount_markup_history
		WHERE cycle_amount_id IN ? AND set_at <= ?
		ORDER BY cycle_amount_id, set_at DESC
	`, cycleAmountIDs, at).Scan(&rows).Error
	result := make(map[int64]types.PricePair, len(rows))
	for _, r := range rows {
		result[r.CycleAmountID] = types.PricePair{RetailPrice: r.RetailMarkup, WholesalePrice: r.WholesaleMarkup}
	}
	return result, err
}
