package repository

import (
	"log"

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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func GetAllOrderProducts(db *gorm.DB, orderID int64) ([]order_products.OrderProduct, error) {
	var products []order_products.OrderProduct
	if err := db.Preload("Product").Where("order_id = ?", orderID).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

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

func GetIndustrialDriveForGate(db *gorm.DB, orderID int64, gateID int64) (*order_gate_industrial_drives.OrderGateIndustrialDrive, error) {
	var drive order_gate_industrial_drives.OrderGateIndustrialDrive
	if err := db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		First(&drive).Error; err != nil {
		return nil, err
	}

	return &drive, nil
}

func GetResidentialDriveForGate(db *gorm.DB, orderID int64, gateID int64) (*order_gate_residential_drives.OrderGateResidentialDrive, error) {
	var drive order_gate_residential_drives.OrderGateResidentialDrive
	if err := db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		First(&drive).Error; err != nil {
		return nil, err
	}

	return &drive, nil
}

func GetManualDriveForGate(db *gorm.DB, orderID int64, gateID int64) (*order_gate_manual_drives.OrderGateManualDrive, error) {
	var drive order_gate_manual_drives.OrderGateManualDrive
	if err := db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		First(&drive).Error; err != nil {
		return nil, err
	}

	return &drive, nil
}

func DeleteOrder(db *gorm.DB, orderID int64) error {
	return db.Where("id = ?", orderID).Delete(&orders.Order{}).Error
}

func CreateNewGate(db *gorm.DB, gate order_gates.OrderGate) (order_gates.OrderGate, error) {
	if err := db.Clauses(clause.Returning{}).Create(&gate).Error; err != nil {
		return order_gates.OrderGate{}, err
	}

	return gate, nil
}

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

func CreateNewOrder(db *gorm.DB, order *orders.Order) error {
	return db.Create(&order).Error
}

func CreateOrderProducts(db *gorm.DB, products []order_products.OrderProduct) error {
	return db.Create(&products).Error
}

func CreateGateOptions(db *gorm.DB, options []order_gate_options.OrderGateOption) error {
	return db.Create(&options).Error
}

func DeleteGateFromOrder(db *gorm.DB, orderID int64, rowNumber int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, rowNumber).
		Delete(&order_gates.OrderGate{}).Error
}

func DeleteGateResidentialDrive(db *gorm.DB, orderID int64, gateID int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(&order_gate_residential_drives.OrderGateResidentialDrive{}).
		Error
}

func DeleteGateIndustrialDrive(db *gorm.DB, orderID int64, gateID int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(order_gate_industrial_drives.OrderGateIndustrialDrive{}).
		Error
}

func DeleteGateManualDrive(db *gorm.DB, orderID int64, gateID int64) error {
	return db.Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(&order_gate_manual_drives.OrderGateManualDrive{}).Error
}

func UpdateGateIndustrialDrive(db *gorm.DB, orderID int64, gateID int64, driveID int64) error {
	return db.Model(&order_gate_industrial_drives.OrderGateIndustrialDrive{}).
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Update("drive_id", driveID).Error
}

func UpdateGateResidentialDrive(db *gorm.DB, orderID int64, gateID int64, driveID int64, railID int64) error {
	return db.Model(&order_gate_residential_drives.OrderGateResidentialDrive{}).
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Updates(map[string]any{"drive_id": driveID, "rail_id": railID}).Error
}

func UpdateGateManualDrive(db *gorm.DB, orderID int64, gateID int64, chainLength int32) error {
	return db.Model(&order_gate_manual_drives.OrderGateManualDrive{}).
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Update("chain_length", chainLength).Error
}

func UpdateGate(db *gorm.DB, gate order_gates.OrderGate) error {
	return db.Save(&gate).Error
}

func UpdateOrderStatus(db *gorm.DB, orderID int64, status enums.OrderStatus) error {
	return db.Model(&orders.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func GetOrderStatus(orderID int64) (*string, error) {
	var status string
	if err := database.DB.Model(&orders.Order{}).
		Where("id = ?", orderID).
		Pluck("status", &status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

func DeleteAllOrderProducts(db *gorm.DB, orderID int64) error {
	return db.
		Where("order_id = ?", orderID).
		Delete(&order_products.OrderProduct{}).Error
}

func GetAllProducts() ([]products.Product, error) {
	var products []products.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func DeleteAllGateOptions(db *gorm.DB, orderID int64, gateID int64) error {
	return db.
		Where("order_id = ? AND row_number = ?", orderID, gateID).
		Delete(&order_gate_options.OrderGateOption{}).Error
}

func GetIndustrialDriveByID(db *gorm.DB, driveID int64) (*industrial_gate_drives.IndustrialGateDrive, error) {
	var industrialDrive industrial_gate_drives.IndustrialGateDrive
	if err := db.Where("id = ?", driveID).First(&industrialDrive).Error; err != nil {
		return nil, err
	}

	return &industrialDrive, nil
}

func GetResidentialDriveByID(db *gorm.DB, driveID int64) (*residential_gate_drives.ResidentialGateDrive, error) {
	var residentialDrive residential_gate_drives.ResidentialGateDrive
	if err := db.Where("id = ?", driveID).First(&residentialDrive).Error; err != nil {
		return nil, err
	}

	return &residentialDrive, nil
}

func GetRailByID(db *gorm.DB, railID int64) (*rails.Rail, error) {
	var rail rails.Rail
	if err := db.Where("id = ?", railID).First(&rail).Error; err != nil {
		return nil, err
	}

	return &rail, nil
}

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

func GetLiftTypeByID(db *gorm.DB, liftTypeID int64) (*lift_types.LiftType, error) {
	var liftType lift_types.LiftType
	if err := db.Where("id = ?", liftTypeID).First(&liftType).Error; err != nil {
		return nil, err
	}

	return &liftType, nil
}

func GetCycleAmountByID(db *gorm.DB, cycleAmountID int64) (*cycle_amounts.CycleAmount, error) {
	var cycleAmount cycle_amounts.CycleAmount
	if err := db.Where("id = ?", cycleAmountID).First(&cycleAmount).Error; err != nil {
		return nil, err
	}

	return &cycleAmount, nil
}

func GetColorByID(db *gorm.DB, colorID int64) (*colors.Color, error) {
	var color colors.Color
	if err := db.Where("id = ?", colorID).First(&color).Error; err != nil {
		return nil, err
	}
	return &color, nil
}

func GetProductByID(db *gorm.DB, productID int64) (*products.Product, error) {
	var product products.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
