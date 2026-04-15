package service

import (
	"errors"
	"project/internal/database"
	"project/internal/domain/enums"
	"project/internal/domain/gates_and_sales_manual_drive"
	"project/internal/domain/gates_and_sales_options"
	"project/internal/domain/industrial_gates_and_sales_drive"
	"project/internal/domain/options"
	"project/internal/domain/residential_gates_and_sales_drive_rail"
	"project/internal/domain/sales"
	"project/internal/domain/sales_and_gates"
	"project/internal/domain/sales_and_products"
	"project/internal/domain/users"
	"project/internal/generated"
	"project/internal/repository"
	"project/internal/types"
)

// Функция для проверки доступа к заказу
func getAccessibleSale(user *users.User, saleID int64) (*sales.Sale, error) {
	var sale sales.Sale
	query := database.DB.Where("id = ?", saleID)

	switch user.Role.Name {
	case "dealer":
		query = query.Where("client_id = ?", user.ID)
	case "manager":
		query = query.Where("manager_id = ?", user.ID)
	default:
		return nil, errors.New("forbidden")
	}

	if err := query.First(&sale).Error; err != nil {
		return nil, err
	}

	return &sale, nil
}

func GetAllUserOrders(user *users.User) ([]sales.Sale, error) {
	// Вернуть limit потом
	if user.Role.Name == "dealer" {
		return repository.GetAllDealerOrders(database.DB, user)
	} else {
		return repository.GetAllManagerOrders(database.DB, user)
	}
}

func GetUserOrderByID(user *users.User, saleID int64) (types.Order, error) {
	_, err := getAccessibleSale(user, saleID)
	if err != nil {
		return types.Order{}, err
	}

	// Получаем все ворота выбранного заказа
	orderGates, err := repository.GetOrderGatesByOrderID(database.DB, saleID)
	if err != nil {
		return types.Order{}, err
	}

	// Собираем все ID ворот
	var gateIDs []int64
	for _, gate := range orderGates {
		gateIDs = append(gateIDs, gate.RowNumber)
	}

	// Получаем все опции всех ворот
	gateOptions, err := repository.GetAllOptionsForOrder(database.DB, gateIDs)
	if err != nil {
		return types.Order{}, err
	}

	// Группируем опции по ID ворот
	optionsMap := make(map[int64][]options.Option)
	for _, opt := range gateOptions {
		optionsMap[opt.RowNumber] = append(optionsMap[opt.RowNumber], opt.Option)
	}

	// Группируем в gates
	order := types.Order{
		Gates:    []types.Gate{},
		Products: []sales_and_products.SalesAndProduct{},
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
		return types.Order{}, err
	}

	return order, nil
}

type CurrentGatePageData struct {
	Gate    sales_and_gates.SalesAndGate
	Options []gates_and_sales_options.GatesAndSalesOption

	IndustrialDrive  *industrial_gates_and_sales_drive.IndustrialGatesAndSalesDrive
	ResidentialDrive *residential_gates_and_sales_drive_rail.ResidentialGatesAndSalesDriveRail
	ManualDrive      *gates_and_sales_manual_drive.GatesAndSalesManualDrive

	Configuraion types.Config
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

	pageData.Configuraion, err = GetGateCfg(pageData.Gate.GateType, true)
	if err != nil {
		return CurrentGatePageData{}, err
	}

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

		gate := sales_and_gates.SalesAndGate{
			SaleID:        saleID,
			GateType:      gateType,
			Width:         int32(cfg.WidthParams.MinValue),
			Height:        int32(cfg.HeightParams.MinValue),
			Headroom:      0,
			LiftTypeID:    cfg.LiftTypes[0].ID,
			CycleAmountID: cfg.CycleAmounts[0].ID,
			ColorOutID:    cfg.Colors[0].ID,
			DriveType:     "industrial",
			Amount:        1,
		}

		gate, err = repository.CreateNewGate(database.DB, gate)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		err = repository.CreateIndustrialDriveForGate(database.DB, saleID, gate.RowNumber, int32(cfg.IndustrialDrives[0].ID))
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		return gate, nil

	case enums.GateTypeRes:
		cfg, err := GetGateCfg(gateType, true)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		gate := sales_and_gates.SalesAndGate{
			SaleID:        saleID,
			GateType:      gateType,
			Width:         int32(cfg.WidthParams.MinValue),
			Height:        int32(cfg.HeightParams.MinValue),
			Headroom:      0,
			LiftTypeID:    cfg.LiftTypes[0].ID,
			CycleAmountID: cfg.CycleAmounts[0].ID,
			ColorOutID:    cfg.Colors[0].ID,
			DriveType:     "residential",
			Amount:        1,
		}

		gate, err = repository.CreateNewGate(database.DB, gate)
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		err = repository.CreateResidentialDriveForGate(database.DB, saleID, gate.RowNumber, int32(cfg.IndustrialDrives[0].ID), int32(cfg.Rails[0].ID))
		if err != nil {
			return sales_and_gates.SalesAndGate{}, err
		}

		return gate, nil

	default:
		return sales_and_gates.SalesAndGate{}, errors.New("неверный тип ворот")
	}
}

func CreateNewOrder(user *users.User, orderData *generated.OrderRequest) error {
	role := user.Role.Name

	tx := database.DB.Begin()

	// Создаем заказ
	var order sales.Sale
	if role == "dealer" {
		managerId, err := repository.GetManagerIdByDealerId(tx, user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		order = sales.Sale{
			ClientID:  &user.ID,
			ManagerID: managerId,
			Status:    enums.OrderStatusNew,
		}
	} else {
		order = sales.Sale{
			ClientID:  nil,
			ManagerID: user.ID,
			Status:    enums.OrderStatusNew,
		}
	}

	if err := repository.CreateNewOrder(tx, &order); err != nil {
		tx.Rollback()
		return err
	}

	// Содаем товары заказа
	if len(orderData.Products) != 0 {
		var salesAndProducts []sales_and_products.SalesAndProduct
		for productId, amount := range orderData.Products {
			salesAndProducts = append(salesAndProducts, sales_and_products.SalesAndProduct{
				SaleID:    order.ID,
				ProductID: productId,
				Amount:    amount,
			})
		}

		if err := repository.CreateOrderProducts(tx, salesAndProducts); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Создаем ворота в заказе
	for i, gate := range orderData.OrderGates {
		gateType, err := enums.GateTypeFromProto(gate.GateType)
		if err != nil {
			tx.Rollback()
			return err
		}

		driveType, err := enums.GetDriveTypeFromProto(gate.Drive)
		if err != nil {
			tx.Rollback()
			return err
		}

		orderDetails := sales_and_gates.SalesAndGate{
			SaleID:        order.ID,
			RowNumber:     int64(i + 1),
			GateType:      gateType,
			Width:         gate.Width,
			Height:        gate.Height,
			Headroom:      gate.Headroom,
			LiftTypeID:    gate.LiftTypeId,
			ColorOutID:    gate.ColorOutId,
			CycleAmountID: gate.CycleAmountId,
			DriveType:     driveType,
			Amount:        gate.Amount,
		}

		newGate, err := repository.CreateNewGate(tx, orderDetails)
		if err != nil {
			tx.Rollback()
			return err
		}

		switch d := gate.Drive.DriveType.(type) {
		case *generated.Drive_Industrial:
			driveID := d.Industrial.DriveId

			if err := repository.CreateIndustrialDriveForGate(tx, order.ID, newGate.RowNumber, int32(driveID)); err != nil {
				tx.Rollback()
				return err
			}

		case *generated.Drive_Residential:
			driveID := d.Residential.DriveId
			railID := d.Residential.RailId

			if err := repository.CreateResidentialDriveForGate(tx, order.ID, newGate.RowNumber, int32(driveID), int32(railID)); err != nil {
				tx.Rollback()
				return err
			}

		case *generated.Drive_Manual:
			chain := d.Manual.ChainLength

			if err := repository.CreateManualDriveForGate(tx, order.ID, newGate.RowNumber, chain); err != nil {
				tx.Rollback()
				return err
			}
		}

		if len(gate.Options) == 0 {
			continue
		}

		// Создаем дополнительные опции у ворот
		var gateAndSalesOptions []gates_and_sales_options.GatesAndSalesOption
		for optionId, amount := range gate.Options {
			gateAndSalesOption := gates_and_sales_options.GatesAndSalesOption{
				SaleID:    order.ID,
				RowNumber: int64(i + 1),
				OptionID:  optionId,
				Amount:    amount,
			}

			gateAndSalesOptions = append(gateAndSalesOptions, gateAndSalesOption)
		}

		if err := repository.CreateGateOptions(tx, gateAndSalesOptions); err != nil {
			tx.Rollback()
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

	gate.Width = gateData.Width
	gate.Height = gateData.Height
	gate.Headroom = gateData.Headroom
	gate.LiftTypeID = gateData.LiftTypeId
	gate.ColorOutID = gateData.ColorOutId
	gate.CycleAmountID = gateData.CycleAmountId
	gate.Amount = gateData.Amount

	switch gateData.Drive.DriveType.(type) {
	case *generated.Drive_Industrial:
		//Если это другой тип привода у ворот
		if gate.DriveType != enums.IndDriveType {
			if err := repository.CreateIndustrialDriveForGate(database.DB, saleID, int64(gateID), int32(gateData.Drive.GetIndustrial().DriveId)); err != nil {
				return err
			}

			// Удаляем предыдущий привод
			switch gate.DriveType {
			case enums.ResDriveType:
				if err := repository.DeleteGateResidentialDrive(database.DB, saleID, int64(gateID)); err != nil {
					return err
				}
			case enums.ManualDriveType:
				if err := repository.DeleteGateManualDrive(database.DB, saleID, int64(gateID)); err != nil {
					return err
				}
			}
			// Если у ворот тот же самый тип привода
		} else {
			if err := repository.UpdateGateIndustrialDrive(database.DB, saleID, int64(gateID), int32(gateData.Drive.GetIndustrial().DriveId)); err != nil {
				return err
			}
		}

	case *generated.Drive_Residential:
		if gate.DriveType != enums.ResDriveType {
			if err := repository.CreateResidentialDriveForGate(database.DB, saleID, int64(gateID), int32(gateData.Drive.GetResidential().DriveId), int32(gateData.Drive.GetResidential().RailId)); err != nil {
				return err
			}

			switch gate.DriveType {
			case enums.IndDriveType:
				if err := repository.DeleteGateIndustrialDrive(database.DB, saleID, int64(gateID)); err != nil {
					return err
				}

			case enums.ManualDriveType:
				if err := repository.DeleteGateManualDrive(database.DB, saleID, int64(gateID)); err != nil {
					return err
				}
			}
		} else {
			if err := repository.UpdateGateResidentialDrive(database.DB, saleID, int64(gateID), int32(gateData.Drive.GetIndustrial().DriveId), int32(gateData.Drive.GetResidential().RailId)); err != nil {
				return err
			}
		}
	case *generated.Drive_Manual:
		if gate.DriveType != enums.ManualDriveType {
			if err := repository.CreateManualDriveForGate(database.DB, saleID, int64(gateID), gateData.Drive.GetManual().ChainLength); err != nil {
				return err
			}

			switch gate.DriveType {
			case enums.ResDriveType:
				if err := repository.DeleteGateResidentialDrive(database.DB, saleID, int64(gateID)); err != nil {
					return err
				}

			case enums.IndDriveType:
				if err := repository.DeleteGateIndustrialDrive(database.DB, saleID, int64(gateID)); err != nil {
					return err
				}
			}
		} else {
			if err := repository.UpdateGateManualDrive(database.DB, saleID, int64(gateID), gateData.Drive.GetManual().ChainLength); err != nil {
				return err
			}
		}
	}

	gate.DriveType, err = enums.GetDriveTypeFromProto(gateData.Drive)

	if err := repository.UpdateGate(database.DB, gate); err != nil {
		return err
	}

	return nil
}
