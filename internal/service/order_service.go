package service

import (
	"errors"
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amount"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/gates_and_sales_manual_drive"
	"github.com/Euhcslel/SagaWeb/internal/domain/gates_and_sales_options"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gates_and_sales_drive"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gates_and_sales_drive_rail"
	"github.com/Euhcslel/SagaWeb/internal/domain/sales"
	"github.com/Euhcslel/SagaWeb/internal/domain/sales_and_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/sales_and_products"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/generated"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Функция для проверки доступа к заказу
func getAccessibleSale(user *users.User, saleID int64) (*sales.Sale, error) {
	var sale sales.Sale
	query := database.DB.Where("id = ?", saleID)

	switch user.Role {
	case enums.DealerRole:
		query = query.Where("client_id = ?", user.ID)
	case enums.ManagerRole, enums.AdminRole, enums.LogisticianRole:
		query = query.Where("manager_id = ?", user.ID)
	default:
		return nil, errs.ErrForbidden
	}

	err := query.First(&sale).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrForbidden
	} else if err != nil {
		return nil, err
	}

	return &sale, nil
}

func GetAllUserOrders(user *users.User) ([]sales.Sale, error) {
	// Вернуть limit потом
	if user.Role == enums.DealerRole {
		return repository.GetAllDealerOrders(database.DB, user)
	} else {
		return repository.GetAllManagerOrders(database.DB, user)
	}
}

type orderPageData struct {
	Order    types.Order
	Products []products.Product
}

func GetOrderPageData(user *users.User, saleID int64) (*orderPageData, error) {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return nil, err
	}

	// Получаем все ворота выбранного заказа
	orderGates, err := repository.GetOrderGatesByOrderID(database.DB, saleID)
	if err != nil {
		return nil, err
	}

	// Собираем все ID ворот
	var gateIDs []int64
	for _, gate := range orderGates {
		gateIDs = append(gateIDs, gate.RowNumber)
	}

	// Получаем все опции всех ворот
	gateOptions, err := repository.GetAllOptionsForOrder(database.DB, saleID, gateIDs)
	if err != nil {
		return nil, err
	}

	// Группируем опции по ID ворот
	optionsMap := make(map[int64][]options.Option)
	for _, opt := range gateOptions {
		optionsMap[opt.RowNumber] = append(optionsMap[opt.RowNumber], opt.Option)
	}

	status, err := repository.GetOrderStatus(user, saleID)
	if err != nil {
		return nil, err
	}

	orderStatus := enums.OrderStatus(status)

	// Группируем в gates
	order := types.Order{
		Gates:    []types.Gate{},
		Products: []sales_and_products.SalesAndProduct{},
		Status:   orderStatus,
	}

	for _, og := range orderGates {
		order.Gates = append(order.Gates, types.Gate{
			Gate:    og,
			Options: optionsMap[og.RowNumber],
		})
	}

	// Получаем все товары
	order.Products, err = repository.GetAllOrderProducts(database.DB, saleID)
	if err != nil {
		return nil, err
	}

	products, err := repository.GetAllProducts()
	if err != nil {
		return nil, err
	}

	pageData := &orderPageData{
		Order:    order,
		Products: products,
	}

	return pageData, nil
}

type CurrentGatePageData struct {
	Gate    sales_and_gates.SalesAndGate
	Options []gates_and_sales_options.GatesAndSalesOption

	IndustrialDrive  *industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive
	ResidentialDrive *residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail
	ManualDrive      *gates_and_sales_manual_drive.GatesAndSalesManualDrive

	Configuration types.Config
	OrderStatus   enums.OrderStatus
}

func GetCurrentGatePageData(user *users.User, saleID int64, gateID int64) (CurrentGatePageData, error) {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	var pageData CurrentGatePageData
	pageData.Gate, err = repository.GetCurrentGate(database.DB, saleID, gateID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	switch pageData.Gate.DriveType {
	case enums.IndDriveType:
		pageData.IndustrialDrive, err = repository.GetIndustrialDriveForGate(database.DB, saleID, gateID)
		if err != nil {
			return CurrentGatePageData{}, err
		}
	case enums.ResDriveType:
		pageData.ResidentialDrive, err = repository.GetResidentialDriveForGate(database.DB, saleID, gateID)
		if err != nil {
			return CurrentGatePageData{}, err
		}
	case enums.ManualDriveType:
		pageData.ManualDrive, err = repository.GetManualDriveForGate(database.DB, saleID, gateID)
		if err != nil {
			return CurrentGatePageData{}, err
		}
	}

	pageData.Options, err = repository.GetCurrentGateOptions(database.DB, saleID, gateID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	pageData.Configuration, err = GetGateCfg(pageData.Gate.GateType, true)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	orderStatus, err := repository.GetOrderStatus(user, saleID)
	if err != nil {
		return CurrentGatePageData{}, err
	}

	pageData.OrderStatus = enums.OrderStatus(orderStatus)

	return pageData, nil
}

func DeleteUserOrder(user *users.User, saleID int64) error {
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	if err := repository.DeleteOrder(database.DB, *sale); err != nil {
		return err
	}

	return nil
}

type PricePair struct {
	RetailPrice    decimal.Decimal
	WholesalePrice decimal.Decimal
}

func AddNewGateInOrder(user *users.User, saleID int64, formGateType string) (sales_and_gates.SalesAndGate, error) {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return sales_and_gates.SalesAndGate{}, err
	}

	gateType := enums.GateType(formGateType)
	switch gateType {
	case enums.GateTypeInd:
		cfg, err := GetGateCfg(gateType, true)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		gatePrice, err := calculateGatePrice(cfg.WidthParams.MinValue, cfg.HeightParams.MinValue,
			gateType, PricePair{WholesalePrice: cfg.IndustrialDrives[0].WholesalePrice,
				RetailPrice: cfg.IndustrialDrives[0].RetailPrice}, cfg.LiftTypes[0], cfg.CycleAmounts[0],
			PricePair{WholesalePrice: decimal.Zero, RetailPrice: decimal.Zero})
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		gate := sales_and_gates.SalesAndGate{
			SaleID:             saleID,
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

		gate, err = repository.CreateNewGate(database.DB, gate)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		err = repository.CreateIndustrialDriveForGate(database.DB, saleID, gate.RowNumber, int64(cfg.IndustrialDrives[0].ID))
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		return gate, nil

	case enums.GateTypeRes:
		cfg, err := GetGateCfg(gateType, true)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		gatePrice, err := calculateGatePrice(cfg.WidthParams.MinValue, cfg.HeightParams.MinValue,
			gateType, PricePair{WholesalePrice: cfg.ResidentialDrives[0].WholesalePrice,
				RetailPrice: cfg.ResidentialDrives[0].RetailPrice}, cfg.LiftTypes[0], cfg.CycleAmounts[0],
			PricePair{WholesalePrice: decimal.Zero, RetailPrice: decimal.Zero})
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		gate := sales_and_gates.SalesAndGate{
			SaleID:             saleID,
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

		gate, err = repository.CreateNewGate(database.DB, gate)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		err = repository.CreateResidentialDriveForGate(database.DB, saleID, gate.RowNumber, cfg.ResidentialDrives[0].ID, cfg.Rails[0].ID)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		return gate, nil

	default:
		return sales_and_gates.SalesAndGate{}, errs.ErrInvalidGateType
	}
}

func calculateGatePrice(width, height int64, gateType enums.GateType,
	drivePrices PricePair, liftType lift_types.LiftType,
	cycleAmount cycle_amount.CycleAmount, optionsPrices PricePair) (*PricePair, error) {
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

	return &PricePair{RetailPrice: gateRetailPrice, WholesalePrice: gateWholesalePrice}, nil
}

func CreateNewOrder(user *users.User, orderData *generated.OrderRequest) error {
	role := user.Role

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Создаем заказ
	var order sales.Sale
	if role == enums.DealerRole {
		managerId, err := repository.GetManagerIdByDealerId(tx, user.ID)
		if err != nil {
			return err
		}

		order = sales.Sale{
			ClientID:  &user.ID,
			ManagerID: managerId,
			Status:    enums.OrderStatusPending,
		}
	} else {
		order = sales.Sale{
			ClientID:  nil,
			ManagerID: user.ID,
			Status:    enums.OrderStatusPending,
		}
	}

	if err := repository.CreateNewOrder(tx, &order); err != nil {
		return err
	}

	// Содаем товары заказа
	if len(orderData.Products) != 0 {
		var salesAndProducts []sales_and_products.SalesAndProduct
		for _, product := range orderData.Products {
			salesAndProducts = append(salesAndProducts, sales_and_products.SalesAndProduct{
				SaleID:    order.ID,
				ProductID: product.ProductId,
				Amount:    product.Amount,
			})
		}

		if err := repository.CreateOrderProducts(tx, salesAndProducts); err != nil {
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

		var drivePrices PricePair
		switch d := gate.Drive.DriveType.(type) {
		case *generated.Drive_Industrial:
			driveID := d.Industrial.DriveId
			drive, err := repository.GetIndustrialDriveById(tx, driveID)
			if err != nil {
				return err
			}

			drivePrices = PricePair{WholesalePrice: drive.WholesalePrice, RetailPrice: drive.RetailPrice}

		case *generated.Drive_Residential:
			driveID := d.Residential.DriveId
			railID := d.Residential.RailId
			drive, err := repository.GetResidentialDriveById(tx, driveID)
			if err != nil {
				return err
			}

			rail, err := repository.GetRailById(tx, railID)
			if err != nil {
				return err
			}

			drivePrices = PricePair{WholesalePrice: drive.WholesalePrice.Add(rail.WholesalePrice),
				RetailPrice: drive.RetailPrice.Add(rail.RetailPrice)}

		case *generated.Drive_Manual:
			chain := d.Manual.ChainLength
			manualPrices, err := repository.GetManualDrivePrices()
			if err != nil {
				return err
			}

			drivePrices = PricePair{
				WholesalePrice: manualPrices.ChainMeterWholesalePrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpWholesalePrice),
				RetailPrice:    manualPrices.ChainMeterRetailPrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpRetailPrice)}
		}

		liftType, err := repository.GetLiftTypeById(tx, gate.LiftTypeId)
		if err != nil {
			return err
		}

		cycleAmount, err := repository.GetCycleAmountById(tx, gate.CycleAmountId)
		if err != nil {
			return err
		}

		var optionsPrices PricePair

		if len(gate.Options) > 0 {
			optionIDs := make([]int64, 0, len(gate.Options))

			for _, gateOption := range gate.Options {
				if gateOption.OptionId == 0 || gateOption.Amount <= 0 {
					continue
				}

				optionIDs = append(optionIDs, gateOption.OptionId)
			}

			if len(optionIDs) > 0 {
				optionsList, err := repository.GetOptionsByIDs(tx, optionIDs)
				if err != nil {
					return err
				}

				optionsByID := make(map[int64]options.Option, len(optionsList))

				for _, option := range optionsList {
					optionsByID[option.ID] = option
				}

				for _, gateOption := range gate.Options {
					if gateOption.OptionId == 0 || gateOption.Amount <= 0 {
						continue
					}

					option, ok := optionsByID[gateOption.OptionId]
					if !ok {
						return errors.New("option not found")
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

		gatePrices, err := calculateGatePrice(int64(gate.Width), int64(gate.Height), gateType,
			drivePrices, liftType, cycleAmount, optionsPrices)
		if err != nil {
			return err
		}

		orderDetails := sales_and_gates.SalesAndGate{
			SaleID:             order.ID,
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
		var gateAndSalesOptions []gates_and_sales_options.GatesAndSalesOption
		for _, gateOption := range gate.Options {
			gateAndSalesOption := gates_and_sales_options.GatesAndSalesOption{
				SaleID:    order.ID,
				RowNumber: int64(i + 1),
				OptionID:  gateOption.OptionId,
				Amount:    gateOption.Amount,
			}

			gateAndSalesOptions = append(gateAndSalesOptions, gateAndSalesOption)
		}

		if err := repository.CreateGateOptions(tx, gateAndSalesOptions); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func DeleteGateFromOrder(user *users.User, saleID int64, rowNumber int64) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	gate, err := repository.GetCurrentGate(database.DB, saleID, rowNumber)
	if err != nil {
		return err
	}

	if err := repository.DeleteGateFromOrder(database.DB, gate); err != nil {
		return err
	}

	return nil
}

func UpdateGateInOrder(user *users.User, saleID int64, gateID int64, gateData *generated.GateConfig) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	gate, err := repository.GetCurrentGate(database.DB, saleID, gateID)
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

	liftType, err := repository.GetLiftTypeById(tx, gateData.LiftTypeId)
	if err != nil {
		return err
	}

	cycleAmount, err := repository.GetCycleAmountById(tx, gateData.CycleAmountId)
	if err != nil {
		return err
	}

	var drivePrices PricePair
	switch gateData.Drive.DriveType.(type) {
	case *generated.Drive_Industrial:
		//Если это другой тип привода у ворот
		if gate.DriveType != enums.IndDriveType {
			if err := repository.CreateIndustrialDriveForGate(tx, saleID, gateID, gateData.Drive.GetIndustrial().DriveId); err != nil {
				return err
			}

			// Удаляем предыдущий привод
			switch gate.DriveType {
			case enums.ResDriveType:
				if err := repository.DeleteGateResidentialDrive(tx, saleID, gateID); err != nil {
					return err
				}
			case enums.ManualDriveType:
				if err := repository.DeleteGateManualDrive(tx, saleID, gateID); err != nil {
					return err
				}
			}
			// Если у ворот тот же самый тип привода
		} else {
			if err := repository.UpdateGateIndustrialDrive(tx, saleID, gateID, gateData.Drive.GetIndustrial().DriveId); err != nil {
				return err
			}
		}

		driveID := gateData.Drive.GetIndustrial().DriveId
		drive, err := repository.GetIndustrialDriveById(tx, driveID)
		if err != nil {
			return err
		}

		drivePrices = PricePair{WholesalePrice: drive.WholesalePrice, RetailPrice: drive.RetailPrice}

	case *generated.Drive_Residential:
		if gate.DriveType != enums.ResDriveType {
			if err := repository.CreateResidentialDriveForGate(tx, saleID, gateID, gateData.Drive.GetResidential().DriveId, gateData.Drive.GetResidential().RailId); err != nil {
				return err
			}

			switch gate.DriveType {
			case enums.IndDriveType:
				if err := repository.DeleteGateIndustrialDrive(tx, saleID, gateID); err != nil {
					return err
				}

			case enums.ManualDriveType:
				if err := repository.DeleteGateManualDrive(tx, saleID, gateID); err != nil {
					return err
				}
			}
		} else {
			if err := repository.UpdateGateResidentialDrive(tx, saleID, gateID, gateData.Drive.GetResidential().DriveId, gateData.Drive.GetResidential().RailId); err != nil {
				return err
			}
		}

		driveID := gateData.Drive.GetResidential().DriveId
		railID := gateData.Drive.GetResidential().RailId
		drive, err := repository.GetResidentialDriveById(tx, driveID)
		if err != nil {
			return err
		}

		rail, err := repository.GetRailById(tx, railID)
		if err != nil {
			return err
		}

		drivePrices = PricePair{WholesalePrice: drive.WholesalePrice.Add(rail.WholesalePrice),
			RetailPrice: drive.RetailPrice.Add(rail.RetailPrice)}

	case *generated.Drive_Manual:
		if gate.DriveType != enums.ManualDriveType {
			if err := repository.CreateManualDriveForGate(tx, saleID, gateID, gateData.Drive.GetManual().ChainLength); err != nil {
				return err
			}

			switch gate.DriveType {
			case enums.ResDriveType:
				if err := repository.DeleteGateResidentialDrive(tx, saleID, gateID); err != nil {
					return err
				}

			case enums.IndDriveType:
				if err := repository.DeleteGateIndustrialDrive(tx, saleID, gateID); err != nil {
					return err
				}
			}
		} else {
			if err := repository.UpdateGateManualDrive(tx, saleID, gateID, gateData.Drive.GetManual().ChainLength); err != nil {
				return err
			}
		}

		chain := gateData.Drive.GetManual().ChainLength
		manualPrices, err := repository.GetManualDrivePrices()
		if err != nil {
			return err
		}

		drivePrices = PricePair{
			WholesalePrice: manualPrices.ChainMeterWholesalePrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpWholesalePrice),
			RetailPrice:    manualPrices.ChainMeterRetailPrice.Mul(decimal.NewFromInt32(chain)).Add(manualPrices.RcpRetailPrice)}
	}

	gate.DriveType, err = enums.GetDriveTypeFromProto(gateData.Drive)
	if err != nil {
		return err
	}

	var optionsList []gates_and_sales_options.GatesAndSalesOption
	for _, option := range gateData.Options {
		optionsList = append(optionsList, gates_and_sales_options.GatesAndSalesOption{
			SaleID:    saleID,
			RowNumber: gateID,
			OptionID:  option.OptionId,
			Amount:    option.Amount,
		})
	}

	var optionsPrices PricePair

	if len(gateData.Options) > 0 {
		optionIDs := make([]int64, 0, len(gateData.Options))

		for _, gateOption := range gateData.Options {
			if gateOption.OptionId == 0 || gateOption.Amount <= 0 {
				continue
			}

			optionIDs = append(optionIDs, gateOption.OptionId)
		}

		if len(optionIDs) > 0 {
			optionsList, err := repository.GetOptionsByIDs(tx, optionIDs)
			if err != nil {
				return err
			}

			optionsByID := make(map[int64]options.Option, len(optionsList))

			for _, option := range optionsList {
				optionsByID[option.ID] = option
			}

			for _, gateOption := range gateData.Options {
				if gateOption.OptionId == 0 || gateOption.Amount <= 0 {
					continue
				}

				option, ok := optionsByID[gateOption.OptionId]
				if !ok {
					return errors.New("option not found")
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

	if err := repository.DeleteAllGateOptions(tx, saleID, gateID); err != nil {
		return err
	}

	if len(optionsList) > 0 {
		if err := repository.CreateGateOptions(tx, optionsList); err != nil {
			return err
		}
	}

	gatePrices, err := calculateGatePrice(int64(gateData.Width), int64(gateData.Height),
		gate.GateType, drivePrices, liftType, cycleAmount, optionsPrices)
	if err != nil {
		return err
	}

	gate.GateRetailPrice = gatePrices.RetailPrice
	gate.GateWholesalePrice = gatePrices.WholesalePrice

	if err := repository.UpdateGate(tx, gate); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func UpdateOrderStatus(user *users.User, saleID int64, updateStatusRequest *generated.UpdateOrderStatusRequest) error {
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	status, err := enums.GetOrderStatusFromProto(&updateStatusRequest.Status)
	if err != nil {
		return err
	}

	if err = repository.UpdateOrderStatus(database.DB, sale, status); err != nil {
		return err
	}

	return nil
}

func UploadOfferToOrder(user *users.User, saleID int64, file multipart.File, handler *multipart.FileHeader) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	offerDirectoryPath := os.Getenv("OFFERS_DIRECTORY")
	if offerDirectoryPath == "" {
		return err
	}

	dst, err := os.Create(offerDirectoryPath + "/" + handler.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	ext := filepath.Ext(handler.Filename)
	name := strings.TrimSuffix(handler.Filename, ext)
	if err = repository.AttachOfferToOrder(database.DB, saleID, name, handler.Filename); err != nil {
		return err
	}

	return nil
}

func UploadContractToOrder(user *users.User, saleID int64, file multipart.File, handler *multipart.FileHeader) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	contractDirectoryPath := os.Getenv("CONTRACTS_DIRECTORY")
	if contractDirectoryPath == "" {
		return err
	}

	dst, err := os.Create(contractDirectoryPath + "/" + handler.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	ext := filepath.Ext(handler.Filename)
	name := strings.TrimSuffix(handler.Filename, ext)
	if err = repository.AttachContractToOrder(database.DB, saleID, name, handler.Filename); err != nil {
		return err
	}

	return nil
}

func UploadBillToOrder(user *users.User, saleID int64, file multipart.File, handler *multipart.FileHeader) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	billDirectoryPath := os.Getenv("BILLS_DIRECTORY")
	if billDirectoryPath == "" {
		return err
	}

	dst, err := os.Create(billDirectoryPath + "/" + handler.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	ext := filepath.Ext(handler.Filename)
	name := strings.TrimSuffix(handler.Filename, ext)
	if err = repository.AttachBillToOrder(database.DB, saleID, name, handler.Filename); err != nil {
		return err
	}

	return nil
}

func GetAllOrderDocuments(user *users.User, saleID int64) (*generated.DocumentsList, error) {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return nil, err
	}

	var resp *generated.DocumentsList

	offersNumberList, err := repository.GetOffersNumberList(database.DB, saleID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	offers := make([]*generated.Offer, 0, len(offersNumberList))
	for _, number := range offersNumberList {
		offers = append(offers, &generated.Offer{
			OfferNumber: number,
		})
	}

	contractsNumberList, err := repository.GetContractsNumberList(database.DB, saleID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	contracts := make([]*generated.Contract, 0, len(contractsNumberList))
	for _, number := range contractsNumberList {
		contracts = append(contracts, &generated.Contract{
			ContractNumber: number,
		})
	}

	billsNumberList, err := repository.GetBillsNumberList(database.DB, saleID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	bills := make([]*generated.Bill, 0, len(billsNumberList))
	for _, number := range billsNumberList {
		bills = append(bills, &generated.Bill{
			BillNumber: number,
		})
	}

	documentsNameList, err := repository.GetDocumentsNameList(database.DB, saleID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	documents := make([]*generated.Document, 0, len(documentsNameList))
	for _, name := range documentsNameList {
		documents = append(documents, &generated.Document{
			Name: name,
		})
	}

	resp = &generated.DocumentsList{
		Documents: documents,
		Offers:    offers,
		Bills:     bills,
		Contracts: contracts,
	}

	return resp, nil
}

type FileInfo struct {
	FilePath string
	FileName string
}

func GetFileInfo(user *users.User, saleID int64, docType string, docName string) (*FileInfo, error) {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return nil, err
	}

	documentType, err := enums.GetDocumentTypeFromString(docType)
	if err != nil {
		return nil, err
	}

	fileInfo := FileInfo{}

	switch documentType {
	case enums.BillDocumentType:
		fileInfo.FileName, err = repository.GetBillFileName(docName)
		if err != nil {
			return nil, err
		}

		billDirectoryPath := os.Getenv("BILLS_DIRECTORY")
		if billDirectoryPath == "" {
			return nil, err
		}

		fileInfo.FilePath = billDirectoryPath + "/" + fileInfo.FileName

		return &fileInfo, nil
	case enums.OfferDocumentType:
		fileInfo.FileName, err = repository.GetOfferFileName(docName)
		if err != nil {
			return nil, err
		}

		offerDirectoryPath := os.Getenv("OFFERS_DIRECTORY")
		if offerDirectoryPath == "" {
			return nil, err
		}
		fileInfo.FilePath = offerDirectoryPath + "/" + fileInfo.FileName

		return &fileInfo, nil
	case enums.ContractDocumentType:
		fileInfo.FileName, err = repository.GetContractFileName(docName)
		if err != nil {
			return nil, err
		}

		contractDirectoryPath := os.Getenv("CONTRACTS_DIRECTORY")
		if contractDirectoryPath == "" {
			return nil, err
		}
		fileInfo.FilePath = contractDirectoryPath + "/" + fileInfo.FileName

		return &fileInfo, nil
	case enums.OtherDocumentType:
		fileInfo.FileName, err = repository.GetDocumentFileName(docName)
		if err != nil {
			return nil, err
		}

		documentsDirectoryPath := os.Getenv("DOCUMENTS_DIRECTORY")
		if documentsDirectoryPath == "" {
			return nil, err
		}

		fileInfo.FilePath = documentsDirectoryPath + "/" + fileInfo.FileName

		return &fileInfo, nil

	default:
		return nil, errs.ErrInvalidDocumentType
	}
}

func DeleteOrderDocument(user *users.User, saleID int64, docType string, docName string) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	documentType, err := enums.GetDocumentTypeFromString(docType)
	if err != nil {
		return err
	}

	switch documentType {
	case enums.BillDocumentType:
		fileName, err := repository.GetBillFileName(docName)
		if err != nil {
			return err
		}

		billDirectoryPath := os.Getenv("BILLS_DIRECTORY")
		if billDirectoryPath == "" {
			return err
		}

		os.Remove(billDirectoryPath + "/" + fileName)

		if err := repository.DeleteOrderBill(docName); err != nil {
			return err
		}
	case enums.OfferDocumentType:
		fileName, err := repository.GetOfferFileName(docName)
		if err != nil {
			return err
		}

		offerDirectoryPath := os.Getenv("OFFERS_DIRECTORY")
		if offerDirectoryPath == "" {
			return err
		}

		os.Remove(offerDirectoryPath + "/" + fileName)

		if err := repository.DeleteOrderOffer(docName); err != nil {
			return err
		}
	case enums.ContractDocumentType:
		fileName, err := repository.GetContractFileName(docName)
		if err != nil {
			return err
		}

		contractDirectoryPath := os.Getenv("CONTRACTS_DIRECTORY")
		if contractDirectoryPath == "" {
			return err
		}

		os.Remove(contractDirectoryPath + "/" + fileName)

		if err := repository.DeleteOrderContract(docName); err != nil {
			return err
		}
	case enums.OtherDocumentType:
		fileName, err := repository.GetDocumentFileName(docName)
		if err != nil {
			return err
		}

		documentsDirectoryPath := os.Getenv("DOCUMENTS_DIRECTORY")
		if documentsDirectoryPath == "" {
			return err
		}

		os.Remove(documentsDirectoryPath + "/" + fileName)

		if err := repository.DeleteOrderDocument(docName); err != nil {
			return err
		}
	}

	return nil
}

func UpdateProductsInOrder(user *users.User, saleID int64, updateProductsRequest *generated.UpdateProductsRequest) error {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return err
	}

	var productList []sales_and_products.SalesAndProduct
	for _, product := range updateProductsRequest.Products {
		productList = append(productList, sales_and_products.SalesAndProduct{
			SaleID:    saleID,
			ProductID: product.ProductId,
			Amount:    product.Amount,
		})
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := repository.DeleteAllOrderProducts(tx, saleID); err != nil {
		return err
	}

	if len(productList) > 0 {
		if err := repository.CreateOrderProducts(tx, productList); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
