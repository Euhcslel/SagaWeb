package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
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
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/generated"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/types"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// checkOrderAccess проверяет доступ к заказу.
func checkOrderAccess(user *users.User, orderID int64) error {
	var order orders.Order
	query := database.DB.Where("id = ?", orderID)

	switch user.Role {
	case enums.DealerRole:
		query = query.Where("dealer_id = ?", user.ID)
	case enums.ManagerRole:
		query = query.Where("manager_id = ?", user.ID)
	case enums.AdminRole, enums.LogisticianRole:
		// логистик видит все заказы без фильтра
	default:
		return errs.ErrForbidden
	}

	err := query.First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.ErrForbidden
	} else if err != nil {
		return err
	}

	return nil
}

// GetAllUserOrders возвращает все заказы пользователя в зависимости от его роли.
func GetAllUserOrders(user *users.User) ([]orders.Order, error) {
	switch user.Role {
	case enums.DealerRole:
		return repository.GetAllDealerOrders(database.DB, user)
	case enums.LogisticianRole, enums.AdminRole:
		return repository.GetAllOrders()
	default:
		return repository.GetAllManagerOrders(database.DB, user)
	}
}

// orderPageData описывает данные необходимые для страницы заказа.
type orderPageData struct {
	Order    types.Order
	Products []products.Product
}

// GetOrderPageData возвращает данные для страницы заказа.
func GetOrderPageData(user *users.User, orderID int64) (*orderPageData, error) {
	if err := checkOrderAccess(user, orderID); err != nil {
		return nil, err
	}

	// Получаем все ворота выбранного заказа
	orderGates, err := repository.GetOrderGatesByOrderID(database.DB, orderID)
	if err != nil {
		return nil, err
	}

	// Собираем все ID ворот
	var gateIDs []int64
	for _, gate := range orderGates {
		gateIDs = append(gateIDs, gate.RowNumber)
	}

	// Получаем все опции всех ворот
	gateOptions, err := repository.GetAllOptionsForOrder(database.DB, orderID, gateIDs)
	if err != nil {
		return nil, err
	}

	// Группируем опции по ID ворот
	optionsMap := make(map[int64][]options.Option)
	for _, opt := range gateOptions {
		optionsMap[opt.RowNumber] = append(optionsMap[opt.RowNumber], opt.Option)
	}

	status, err := repository.GetOrderStatus(orderID)
	if err != nil {
		return nil, err
	}

	orderStatus := enums.OrderStatus(*status)

	manufactureDate, err := repository.GetOrderManufactureDate(database.DB, orderID)
	if err != nil {
		return nil, err
	}

	// Группируем в gates
	order := types.Order{
		Gates:           []types.Gate{},
		Products:        []order_products.OrderProduct{},
		Status:          orderStatus,
		ManufactureDate: manufactureDate,
	}

	for _, og := range orderGates {
		order.Gates = append(order.Gates, types.Gate{
			Gate:    og,
			Options: optionsMap[og.RowNumber],
		})
	}

	// Получаем все товары
	order.Products, err = repository.GetAllOrderProducts(database.DB, orderID)
	if err != nil {
		return nil, err
	}

	products, err := repository.GetAllProducts()
	if err != nil {
		return nil, err
	}

	// Для неактивных заказов подгружаем исторические цены товаров
	if orderStatus != enums.OrderStatusPending {
		finalizedAt, err := repository.GetOrderFinalizedAt(database.DB, orderID)
		if err != nil {
			return nil, err
		}

		if finalizedAt != nil {
			productIDs := make([]int64, len(products))
			for i, p := range products {
				productIDs[i] = p.ID
			}

			historicalPrices, err := repository.GetProductHistoricalPrices(database.DB, productIDs, *finalizedAt)
			if err != nil {
				return nil, err
			}

			for i, p := range products {
				if prices, ok := historicalPrices[p.ID]; ok {
					products[i].RetailPrice = prices.RetailPrice
					products[i].WholesalePrice = prices.WholesalePrice
				}
			}
		}
	}

	pageData := &orderPageData{
		Order:    order,
		Products: products,
	}

	return pageData, nil
}

// CurrentGatePageData описывает данные необходимые для страницы ворот в заказе.
type CurrentGatePageData struct {
	Gate    *order_gates.OrderGate
	Options []order_gate_options.OrderGateOption

	IndustrialDrive  *order_gate_industrial_drives.OrderGateIndustrialDrive
	ResidentialDrive *order_gate_residential_drives.OrderGateResidentialDrive
	ManualDrive      *order_gate_manual_drives.OrderGateManualDrive

	Configuration *types.Config
	OrderStatus   enums.OrderStatus
}

// GetCurrentGatePageData возвращает данные для страницы ворот в заказе.
func GetCurrentGatePageData(user *users.User, orderID int64, gateID int64) (CurrentGatePageData, error) {
	err := checkOrderAccess(user, orderID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	var pageData CurrentGatePageData
	pageData.Gate, err = repository.GetCurrentGate(database.DB, orderID, gateID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	switch pageData.Gate.DriveType {
	case enums.IndDriveType:
		pageData.IndustrialDrive, err = repository.GetIndustrialDriveForGate(database.DB, orderID, gateID)
		if err != nil {
			return CurrentGatePageData{}, err
		}
	case enums.ResDriveType:
		pageData.ResidentialDrive, err = repository.GetResidentialDriveForGate(database.DB, orderID, gateID)
		if err != nil {
			return CurrentGatePageData{}, err
		}
	case enums.ManualDriveType:
		pageData.ManualDrive, err = repository.GetManualDriveForGate(database.DB, orderID, gateID)
		if err != nil {
			return CurrentGatePageData{}, err
		}
	}

	pageData.Options, err = repository.GetCurrentGateOptions(database.DB, orderID, gateID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	pageData.Configuration, err = GetGateCfg(pageData.Gate.GateType, true)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	orderStatus, err := repository.GetOrderStatus(orderID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	pageData.OrderStatus = enums.OrderStatus(*orderStatus)

	// Для неактивных заказов подгружаем исторические цены из истории
	if pageData.OrderStatus != enums.OrderStatusPending {
		finalizedAt, err := repository.GetOrderFinalizedAt(database.DB, orderID)
		if err != nil {
			return CurrentGatePageData{}, err
		}

		if finalizedAt != nil {
			if err = applyHistoricalConfigPrices(database.DB, pageData.Configuration, *finalizedAt); err != nil {
				return CurrentGatePageData{}, err
			}
		}
	}

	return pageData, nil
}

// applyHistoricalConfigPrices обновляет цены в конфигурации на основе исторических данных.
func applyHistoricalConfigPrices(db *gorm.DB, cfg *types.Config, at time.Time) error {
	// Опции
	if len(cfg.Options) > 0 {
		ids := make([]int64, len(cfg.Options))
		for i, o := range cfg.Options {
			ids[i] = o.ID
		}
		prices, err := repository.GetOptionHistoricalPrices(db, ids, at)
		if err != nil {
			return err
		}
		for i, o := range cfg.Options {
			if p, ok := prices[o.ID]; ok {
				cfg.Options[i].RetailPrice = p.RetailPrice
				cfg.Options[i].WholesalePrice = p.WholesalePrice
			}
		}
	}

	// Промышленные приводы
	if len(cfg.IndustrialDrives) > 0 {
		ids := make([]int64, len(cfg.IndustrialDrives))
		for i, d := range cfg.IndustrialDrives {
			ids[i] = d.ID
		}
		prices, err := repository.GetIndustrialDriveHistoricalPrices(db, ids, at)
		if err != nil {
			return err
		}
		for i, d := range cfg.IndustrialDrives {
			if p, ok := prices[d.ID]; ok {
				cfg.IndustrialDrives[i].RetailPrice = p.RetailPrice
				cfg.IndustrialDrives[i].WholesalePrice = p.WholesalePrice
			}
		}
	}

	// Бытовые приводы
	if len(cfg.ResidentialDrives) > 0 {
		ids := make([]int64, len(cfg.ResidentialDrives))
		for i, d := range cfg.ResidentialDrives {
			ids[i] = d.ID
		}
		prices, err := repository.GetResidentialDriveHistoricalPrices(db, ids, at)
		if err != nil {
			return err
		}
		for i, d := range cfg.ResidentialDrives {
			if p, ok := prices[d.ID]; ok {
				cfg.ResidentialDrives[i].RetailPrice = p.RetailPrice
				cfg.ResidentialDrives[i].WholesalePrice = p.WholesalePrice
			}
		}
	}

	// Направляющие
	if len(cfg.Rails) > 0 {
		ids := make([]int64, len(cfg.Rails))
		for i, r := range cfg.Rails {
			ids[i] = r.ID
		}
		prices, err := repository.GetRailHistoricalPrices(db, ids, at)
		if err != nil {
			return err
		}
		for i, r := range cfg.Rails {
			if p, ok := prices[r.ID]; ok {
				cfg.Rails[i].RetailPrice = p.RetailPrice
				cfg.Rails[i].WholesalePrice = p.WholesalePrice
			}
		}
	}

	// Типы подъёма (наценки)
	if len(cfg.LiftTypes) > 0 {
		ids := make([]int64, len(cfg.LiftTypes))
		for i, l := range cfg.LiftTypes {
			ids[i] = l.ID
		}
		markups, err := repository.GetLiftTypeHistoricalMarkups(db, ids, at)
		if err != nil {
			return err
		}
		for i, l := range cfg.LiftTypes {
			if p, ok := markups[l.ID]; ok {
				cfg.LiftTypes[i].RetailMarkup = p.RetailPrice
				cfg.LiftTypes[i].WholesaleMarkup = p.WholesalePrice
			}
		}
	}

	// Количество циклов (наценки)
	if len(cfg.CycleAmounts) > 0 {
		ids := make([]int64, len(cfg.CycleAmounts))
		for i, c := range cfg.CycleAmounts {
			ids[i] = c.ID
		}
		markups, err := repository.GetCycleAmountHistoricalMarkups(db, ids, at)
		if err != nil {
			return err
		}
		for i, c := range cfg.CycleAmounts {
			if p, ok := markups[c.ID]; ok {
				cfg.CycleAmounts[i].RetailMarkup = p.RetailPrice
				cfg.CycleAmounts[i].WholesaleMarkup = p.WholesalePrice
			}
		}
	}

	return nil
}

// DeleteUserOrder удаляет заказ пользователя.
func DeleteUserOrder(user *users.User, orderID int64) error {
	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}

	if err := repository.DeleteOrder(database.DB, orderID); err != nil {
		return err
	}

	return nil
}

// AddNewGateInOrder добавляет новые ворота в заказ.
func AddNewGateInOrder(user *users.User, orderID int64, formGateType string) (*order_gates.OrderGate, error) {
	if err := checkOrderAccess(user, orderID); err != nil {
		return nil, err
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	gateType := enums.GateType(formGateType)
	switch gateType {
	case enums.GateTypeInd:
		cfg, err := GetGateCfg(gateType, true)
		if err != nil {
			return nil, err
		}

		if len(cfg.IndustrialDrives) == 0 || len(cfg.LiftTypes) == 0 || len(cfg.CycleAmounts) == 0 || len(cfg.Colors) == 0 {
			return nil, fmt.Errorf("конфигурация неполная: в базе данных отсутствуют необходимые справочники")
		}

		gatePrice, err := calculateGatePrice(cfg.WidthParams.MinValue, cfg.HeightParams.MinValue,
			gateType, types.PricePair{WholesalePrice: cfg.IndustrialDrives[0].WholesalePrice,
				RetailPrice: cfg.IndustrialDrives[0].RetailPrice}, cfg.LiftTypes[0], cfg.CycleAmounts[0],
			types.PricePair{WholesalePrice: decimal.Zero, RetailPrice: decimal.Zero})
		if err != nil {
			return nil, err
		}

		gate := order_gates.OrderGate{
			OrderID:            orderID,
			GateType:           gateType,
			Width:              int32(cfg.WidthParams.MinValue),
			Height:             int32(cfg.HeightParams.MinValue),
			Headroom:           0,
			LiftTypeID:         cfg.LiftTypes[0].ID,
			CycleAmountID:      cfg.CycleAmounts[0].ID,
			ColorOutID:         cfg.Colors[0].ID,
			DriveType:          enums.IndDriveType,
			Amount:             1,
			GateRetailPrice:    gatePrice.RetailPrice,
			GateWholesalePrice: gatePrice.WholesalePrice,
		}

		gate, err = repository.CreateNewGate(tx, gate)
		if err != nil {
			return nil, err
		}

		err = repository.CreateIndustrialDriveForGate(tx, orderID, gate.RowNumber, int64(cfg.IndustrialDrives[0].ID))
		if err != nil {
			return nil, err
		}

		return &gate, tx.Commit().Error

	case enums.GateTypeRes:
		cfg, err := GetGateCfg(gateType, true)
		if err != nil {
			return nil, err
		}

		if len(cfg.ResidentialDrives) == 0 || len(cfg.LiftTypes) == 0 || len(cfg.CycleAmounts) == 0 || len(cfg.Colors) == 0 || len(cfg.Rails) == 0 {
			return nil, fmt.Errorf("конфигурация неполная: в базе данных отсутствуют необходимые справочники")
		}

		gatePrice, err := calculateGatePrice(cfg.WidthParams.MinValue, cfg.HeightParams.MinValue,
			gateType, types.PricePair{WholesalePrice: cfg.ResidentialDrives[0].WholesalePrice,
				RetailPrice: cfg.ResidentialDrives[0].RetailPrice}, cfg.LiftTypes[0], cfg.CycleAmounts[0],
			types.PricePair{WholesalePrice: decimal.Zero, RetailPrice: decimal.Zero})
		if err != nil {
			return nil, err
		}

		gate := order_gates.OrderGate{
			OrderID:            orderID,
			GateType:           gateType,
			Width:              int32(cfg.WidthParams.MinValue),
			Height:             int32(cfg.HeightParams.MinValue),
			Headroom:           0,
			LiftTypeID:         cfg.LiftTypes[0].ID,
			CycleAmountID:      cfg.CycleAmounts[0].ID,
			ColorOutID:         cfg.Colors[0].ID,
			DriveType:          enums.ResDriveType,
			Amount:             1,
			GateRetailPrice:    gatePrice.RetailPrice,
			GateWholesalePrice: gatePrice.WholesalePrice,
		}

		gate, err = repository.CreateNewGate(tx, gate)
		if err != nil {
			return nil, err
		}

		err = repository.CreateResidentialDriveForGate(tx, orderID, gate.RowNumber, cfg.ResidentialDrives[0].ID, cfg.Rails[0].ID)
		if err != nil {
			return nil, err
		}

		return &gate, tx.Commit().Error

	default:
		return nil, errs.ErrInvalidGateType
	}
}

// calculateGatePrice рассчитывает цену ворот на основе заданных параметров.
func calculateGatePrice(width, height int64, gateType enums.GateType,
	drivePrices types.PricePair, liftType lift_types.LiftType,
	cycleAmount cycle_amounts.CycleAmount, optionsPrices types.PricePair) (*types.PricePair, error) {
	var gateRetailPrice decimal.Decimal
	var gateWholesalePrice decimal.Decimal

	sizePrice, err := repository.GetSizeForDimensions(width, height, gateType)
	if err != nil {
		return nil, err
	}

	gateRetailPrice = gateRetailPrice.Add(sizePrice.RetailPrice).
		Add(drivePrices.RetailPrice).
		Add(liftType.RetailMarkup.Div(decimal.NewFromInt(100)).Mul(sizePrice.RetailPrice)).
		Add(cycleAmount.RetailMarkup.Div(decimal.NewFromInt(100)).Mul(sizePrice.RetailPrice)).
		Add(optionsPrices.RetailPrice)

	gateWholesalePrice = gateWholesalePrice.Add(sizePrice.WholesalePrice).
		Add(drivePrices.WholesalePrice).
		Add(liftType.WholesaleMarkup.Div(decimal.NewFromInt(100)).Mul(sizePrice.WholesalePrice)).
		Add(cycleAmount.WholesaleMarkup.Div(decimal.NewFromInt(100)).Mul(sizePrice.WholesalePrice)).
		Add(optionsPrices.WholesalePrice)

	return &types.PricePair{RetailPrice: gateRetailPrice, WholesalePrice: gateWholesalePrice}, nil
}

func calculateOptionsPrices(tx *gorm.DB, gateOptions []*generated.Option) (*types.PricePair, error) {
	if len(gateOptions) == 0 {
		return &types.PricePair{}, nil
	}

	var optionsPrices types.PricePair
	if len(gateOptions) > 0 {
		optionIDs := make([]int64, 0, len(gateOptions))

		for _, gateOption := range gateOptions {
			if gateOption.OptionId == 0 || gateOption.Amount <= 0 {
				continue
			}

			optionIDs = append(optionIDs, gateOption.OptionId)
		}

		if len(optionIDs) > 0 {
			optionsList, err := repository.GetOptionsByIDs(tx, optionIDs)
			if err != nil {
				return nil, err
			}

			optionsByID := make(map[int64]options.Option, len(optionsList))

			for _, option := range optionsList {
				optionsByID[option.ID] = option
			}

			for _, gateOption := range gateOptions {
				if gateOption.OptionId == 0 || gateOption.Amount <= 0 {
					continue
				}

				option, ok := optionsByID[gateOption.OptionId]
				if !ok {
					return nil, errors.New("option not found")
				}

				amount := decimal.NewFromInt32(gateOption.Amount)

				optionsPrices.RetailPrice = optionsPrices.RetailPrice.Add(
					option.RetailPrice.Mul(amount),
				)

				optionsPrices.WholesalePrice = optionsPrices.WholesalePrice.Add(
					option.WholesalePrice.Mul(amount),
				)
			}
		}
	}

	return &optionsPrices, nil
}

// CreateNewOrder создает новый заказ.
func CreateNewOrder(user *users.User, orderData *generated.OrderRequest) error {
	if user.Role != enums.DealerRole {
		return errs.ErrForbidden
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Создаем заказ
	managerID, err := repository.GetManagerIDByDealerID(tx, user.ID)
	if err != nil {
		return err
	}

	order := orders.Order{
		DealerID:  &user.ID,
		ManagerID: managerID,
		Status:    enums.OrderStatusPending,
	}

	if err := repository.CreateNewOrder(tx, &order); err != nil {
		return err
	}

	// Содаем товары заказа
	if len(orderData.Products) != 0 {
		var orderProducts []order_products.OrderProduct
		for _, product := range orderData.Products {
			orderProducts = append(orderProducts, order_products.OrderProduct{
				OrderID:   order.ID,
				ProductID: product.ProductId,
				Amount:    product.Amount,
			})
		}

		if err := repository.CreateOrderProducts(tx, orderProducts); err != nil {
			return err
		}
	}

	// Создаем ворота в заказе
	for i, gate := range orderData.OrderGates {
		gateType, err := enums.GateTypeFromProto(gate.GateType)
		if err != nil {
			return err
		}

		driveType, err := enums.GetDriveTypeFromProto(gate.Drive)
		if err != nil {
			return err
		}

		var drivePrices types.PricePair
		switch d := gate.Drive.DriveType.(type) {
		case *generated.Drive_Industrial:
			driveID := d.Industrial.DriveId
			drive, err := repository.GetIndustrialDriveByID(tx, driveID)
			if err != nil {
				return err
			}

			drivePrices = types.PricePair{WholesalePrice: drive.WholesalePrice, RetailPrice: drive.RetailPrice}

		case *generated.Drive_Residential:
			driveID := d.Residential.DriveId
			railID := d.Residential.RailId
			drive, err := repository.GetResidentialDriveByID(tx, driveID)
			if err != nil {
				return err
			}

			rail, err := repository.GetRailByID(tx, railID)
			if err != nil {
				return err
			}

			drivePrices = types.PricePair{WholesalePrice: drive.WholesalePrice.Add(rail.WholesalePrice),
				RetailPrice: drive.RetailPrice.Add(rail.RetailPrice)}

		case *generated.Drive_Manual:
			chain := d.Manual.ChainLength
			manualPrices, err := repository.GetManualDrivePrices()
			if err != nil {
				return err
			}

			drivePrices = types.PricePair{
				WholesalePrice: manualPrices.ChainMeterWholesalePrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpWholesalePrice),
				RetailPrice:    manualPrices.ChainMeterRetailPrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpRetailPrice)}
		}

		liftType, err := repository.GetLiftTypeByID(tx, gate.LiftTypeId)
		if err != nil {
			return err
		}

		cycleAmount, err := repository.GetCycleAmountByID(tx, gate.CycleAmountId)
		if err != nil {
			return err
		}

		optionsPrices, err := calculateOptionsPrices(tx, gate.Options)
		if err != nil {
			return err
		}

		gatePrices, err := calculateGatePrice(int64(gate.Width), int64(gate.Height), gateType,
			drivePrices, *liftType, *cycleAmount, *optionsPrices)
		if err != nil {
			return err
		}

		orderDetails := order_gates.OrderGate{
			OrderID:            order.ID,
			RowNumber:          int64(i + 1),
			GateType:           gateType,
			Width:              gate.Width,
			Height:             gate.Height,
			Headroom:           gate.Headroom,
			LiftTypeID:         gate.LiftTypeId,
			ColorOutID:         gate.ColorOutId,
			CycleAmountID:      gate.CycleAmountId,
			DriveType:          driveType,
			Amount:             gate.Amount,
			GateRetailPrice:    gatePrices.RetailPrice,
			GateWholesalePrice: gatePrices.WholesalePrice,
		}

		newGate, err := repository.CreateNewGate(tx, orderDetails)
		if err != nil {
			return err
		}

		switch d := gate.Drive.DriveType.(type) {
		case *generated.Drive_Industrial:
			driveID := d.Industrial.DriveId

			if err := repository.CreateIndustrialDriveForGate(tx, order.ID, newGate.RowNumber, driveID); err != nil {
				return err
			}

		case *generated.Drive_Residential:
			driveID := d.Residential.DriveId
			railID := d.Residential.RailId

			if err := repository.CreateResidentialDriveForGate(tx, order.ID, newGate.RowNumber, driveID, railID); err != nil {
				return err
			}

		case *generated.Drive_Manual:
			chain := d.Manual.ChainLength

			if err := repository.CreateManualDriveForGate(tx, order.ID, newGate.RowNumber, chain); err != nil {
				return err
			}
		}

		if len(gate.Options) == 0 {
			continue
		}

		// Создаем дополнительные опции у ворот
		var orderGateOptions []order_gate_options.OrderGateOption
		for _, gateOption := range gate.Options {
			orderGateOption := order_gate_options.OrderGateOption{
				OrderID:   order.ID,
				RowNumber: int64(i + 1),
				OptionID:  gateOption.OptionId,
				Amount:    gateOption.Amount,
			}

			orderGateOptions = append(orderGateOptions, orderGateOption)
		}

		if err := repository.CreateGateOptions(tx, orderGateOptions); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteGateFromOrder удаляет ворота из заказа.
func DeleteGateFromOrder(user *users.User, orderID int64, rowNumber int64) error {
	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}

	if err := repository.DeleteGateFromOrder(database.DB, orderID, rowNumber); err != nil {
		return err
	}

	return nil
}

// UpdateGateInOrder обновляет ворота в заказе.
func UpdateGateInOrder(user *users.User, orderID int64, gateID int64, gateData *generated.GateConfig) error {
	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}

	gate, err := repository.GetCurrentGate(database.DB, orderID, gateID)
	if err != nil {
		return err
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	gate.Width = gateData.Width
	gate.Height = gateData.Height
	gate.Headroom = gateData.Headroom
	gate.LiftTypeID = gateData.LiftTypeId
	gate.ColorOutID = gateData.ColorOutId
	gate.CycleAmountID = gateData.CycleAmountId
	gate.Amount = gateData.Amount

	liftType, err := repository.GetLiftTypeByID(tx, gateData.LiftTypeId)
	if err != nil {
		return err
	}

	cycleAmount, err := repository.GetCycleAmountByID(tx, gateData.CycleAmountId)
	if err != nil {
		return err
	}

	var drivePrices types.PricePair
	switch gateData.Drive.DriveType.(type) {
	case *generated.Drive_Industrial:
		//Если это другой тип привода у ворот
		if gate.DriveType != enums.IndDriveType {
			if err := repository.CreateIndustrialDriveForGate(tx, orderID, gateID, gateData.Drive.GetIndustrial().DriveId); err != nil {
				return err
			}

			// Удаляем предыдущий привод
			switch gate.DriveType {
			case enums.ResDriveType:
				if err := repository.DeleteGateResidentialDrive(tx, orderID, gateID); err != nil {
					return err
				}
			case enums.ManualDriveType:
				if err := repository.DeleteGateManualDrive(tx, orderID, gateID); err != nil {
					return err
				}
			}
			// Если у ворот тот же самый тип привода
		} else {
			if err := repository.UpdateGateIndustrialDrive(tx, orderID, gateID, gateData.Drive.GetIndustrial().DriveId); err != nil {
				return err
			}
		}

		driveID := gateData.Drive.GetIndustrial().DriveId
		drive, err := repository.GetIndustrialDriveByID(tx, driveID)
		if err != nil {
			return err
		}

		drivePrices = types.PricePair{WholesalePrice: drive.WholesalePrice, RetailPrice: drive.RetailPrice}

	case *generated.Drive_Residential:
		if gate.DriveType != enums.ResDriveType {
			if err := repository.CreateResidentialDriveForGate(tx, orderID, gateID, gateData.Drive.GetResidential().DriveId, gateData.Drive.GetResidential().RailId); err != nil {
				return err
			}

			switch gate.DriveType {
			case enums.IndDriveType:
				if err := repository.DeleteGateIndustrialDrive(tx, orderID, gateID); err != nil {
					return err
				}

			case enums.ManualDriveType:
				if err := repository.DeleteGateManualDrive(tx, orderID, gateID); err != nil {
					return err
				}
			}
		} else {
			if err := repository.UpdateGateResidentialDrive(tx, orderID, gateID, gateData.Drive.GetResidential().DriveId, gateData.Drive.GetResidential().RailId); err != nil {
				return err
			}
		}

		driveID := gateData.Drive.GetResidential().DriveId
		railID := gateData.Drive.GetResidential().RailId
		drive, err := repository.GetResidentialDriveByID(tx, driveID)
		if err != nil {
			return err
		}

		rail, err := repository.GetRailByID(tx, railID)
		if err != nil {
			return err
		}

		drivePrices = types.PricePair{WholesalePrice: drive.WholesalePrice.Add(rail.WholesalePrice),
			RetailPrice: drive.RetailPrice.Add(rail.RetailPrice)}

	case *generated.Drive_Manual:
		if gate.DriveType != enums.ManualDriveType {
			if err := repository.CreateManualDriveForGate(tx, orderID, gateID, gateData.Drive.GetManual().ChainLength); err != nil {
				return err
			}

			switch gate.DriveType {
			case enums.ResDriveType:
				if err := repository.DeleteGateResidentialDrive(tx, orderID, gateID); err != nil {
					return err
				}

			case enums.IndDriveType:
				if err := repository.DeleteGateIndustrialDrive(tx, orderID, gateID); err != nil {
					return err
				}
			}
		} else {
			if err := repository.UpdateGateManualDrive(tx, orderID, gateID, gateData.Drive.GetManual().ChainLength); err != nil {
				return err
			}
		}

		chain := gateData.Drive.GetManual().ChainLength
		manualPrices, err := repository.GetManualDrivePrices()
		if err != nil {
			return err
		}

		drivePrices = types.PricePair{
			WholesalePrice: manualPrices.ChainMeterWholesalePrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpWholesalePrice),
			RetailPrice:    manualPrices.ChainMeterRetailPrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpRetailPrice)}
	}

	gate.DriveType, err = enums.GetDriveTypeFromProto(gateData.Drive)
	if err != nil {
		return err
	}

	var optionsList []order_gate_options.OrderGateOption
	for _, option := range gateData.Options {
		optionsList = append(optionsList, order_gate_options.OrderGateOption{
			OrderID:   orderID,
			RowNumber: gateID,
			OptionID:  option.OptionId,
			Amount:    option.Amount,
		})
	}

	if err := repository.DeleteAllGateOptions(tx, orderID, gateID); err != nil {
		return err
	}

	if len(optionsList) > 0 {
		if err := repository.CreateGateOptions(tx, optionsList); err != nil {
			return err
		}
	}

	optionsPrices, err := calculateOptionsPrices(tx, gateData.Options)
	if err != nil {
		return err
	}

	gatePrices, err := calculateGatePrice(int64(gateData.Width), int64(gateData.Height),
		gate.GateType, drivePrices, *liftType, *cycleAmount, *optionsPrices)
	if err != nil {
		return err
	}

	gate.GateRetailPrice = gatePrices.RetailPrice
	gate.GateWholesalePrice = gatePrices.WholesalePrice

	if err := repository.UpdateGate(tx, *gate); err != nil {
		return err
	}

	return tx.Commit().Error
}

// UpdateOrderStatus обновляет статус заказа.
func UpdateOrderStatus(user *users.User, orderID int64, updateStatusRequest *generated.UpdateOrderStatusRequest) error {
	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}
	if user.Role == enums.DealerRole {
		return errs.ErrForbidden
	}

	status, err := enums.GetOrderStatusFromProto(&updateStatusRequest.Status)
	if err != nil {
		return err
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err = repository.UpdateOrderStatus(tx, orderID, status); err != nil {
		return err
	}

	if status != enums.OrderStatusPending {
		if err = repository.SetOrderFinalizedAt(tx, orderID); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

// UpdateProductsInOrder обновляет список товаров в заказе.
func UpdateProductsInOrder(user *users.User, orderID int64, updateProductsRequest *generated.UpdateProductsRequest) error {
	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}

	var productList []order_products.OrderProduct
	for _, product := range updateProductsRequest.Products {
		productList = append(productList, order_products.OrderProduct{
			OrderID:   orderID,
			ProductID: product.ProductId,
			Amount:    product.Amount,
		})
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := repository.DeleteAllOrderProducts(tx, orderID); err != nil {
		return err
	}

	if len(productList) > 0 {
		if err := repository.CreateOrderProducts(tx, productList); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
