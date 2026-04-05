package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"project/pkg/proto_files"
	"project/pkg/types"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm/clause"
)

// Функция для проверки доступа к заказу
func getAccessibleSale(user *models.User, saleID string) (*models.Sale, error) {
	var sale models.Sale
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

// Route: /orders
// Method: GET
func GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	// Вернуть limit потом
	var orders []models.Sale
	if user.Role.Name == "dealer" {
		if err := database.DB.
			Preload("Client").
			Preload("Manager").
			Where("client_id = ?", user.ID).
			Find(&orders).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		if err := database.DB.
			Preload("Client").
			Preload("Manager").
			Where("manager_id = ?", user.ID).
			Find(&orders).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}

	data := map[string]any{
		"css":    "order_list.css",
		"orders": orders,
		"user":   user,
	}

	if err := templates.ExecuteTemplate(w, "list.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /api/orders
// Method: GET
func GetAllUserOrdersAPI(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: GET
func GetUserOrderById(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	saleID := mux.Vars(r)["order_id"]
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	}
	// Получаем все ворота выбранного заказа
	var orderGates []models.SalesAndGate
	if err := database.DB.
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("sale_id = ?", sale.ID).
		Find(&orderGates).Error; err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	// Собираем все ID ворот
	var gateIDs []int64
	for _, gate := range orderGates {
		gateIDs = append(gateIDs, gate.RowNumber)
	}

	// Получаем все опции всех ворот
	var options []models.GatesAndSalesOption
	if err := database.DB.
		Where("row_number IN ?", gateIDs).
		Preload("Option").
		Find(&options).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	// Группируем опции по ID ворот
	optionsMap := make(map[int64][]models.Option)
	for _, opt := range options {
		optionsMap[opt.RowNumber] = append(optionsMap[opt.RowNumber], opt.Option)
	}

	// Группируем в gates
	order := types.Order{
		Gates:    []types.Gate{},
		Products: []models.SalesAndProduct{},
	}
	for _, og := range orderGates {
		order.Gates = append(order.Gates, types.Gate{
			Gate:    og,
			Options: optionsMap[og.RowNumber],
		})
	}

	if err := database.DB.Preload("Product").Where("sale_id = ?", sale.ID).Find(&order.Products).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":     "/order.css",
		"order":   order,
		"user":    user,
		"orderId": sale.ID,
	}

	if err := templates.ExecuteTemplate(w, "order.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

func isManager(role string) bool {
	return role == "manager"
}

// Route: /orders/{order_id}/{gate_id}
// Method: GET
func GetGateInOrder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	saleID := vars["order_id"]
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	}

	var gate models.SalesAndGate
	if err := database.DB.Model(models.SalesAndGate{}).
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("sale_id = ? and row_number = ?", sale.ID, vars["gate_id"]).
		First(&gate).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var drive any
	switch gate.DriveType {
	case models.IndDriveType:
		var d models.IndustrialGatesAndSalesDrive
		if err = database.DB.
			Where("sale_id = ? AND row_number = ?", gate.SaleID, gate.RowNumber).
			First(&d).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		drive = d
	case models.ResDriveType:
		var d models.ResidentialGatesAndSalesDriveRail
		if err = database.DB.
			Where("sale_id = ? AND row_number = ?", gate.SaleID, gate.RowNumber).
			First(&d).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		drive = d
	case models.ManualDriveType:
		var d models.GatesAndSalesManualDrive
		if err = database.DB.
			Where("sale_id = ? AND row_number = ?", gate.SaleID, gate.RowNumber).
			First(&d).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		drive = d
	}

	var options []models.GatesAndSalesOption
	if err := database.DB.Model(models.GatesAndSalesOption{}).Preload("Option").Where("sale_id = ? and row_number = ?", sale.ID, vars["gate_id"]).Find(&options).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	cfg, err := getGateCfg(gate.GateType, true)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":      "calc.css",
		"gate":     gate,
		"options":  options,
		"cfg":      cfg,
		"drive":    drive,
		"statuses": models.GetAllOrderStatuses(),
		"user":     user,
	}

	if err := templates.ExecuteTemplate(w, "gate.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/documents
// Method: GET
func GetOrderDocuments(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: DELETE
func DeleteUserOrder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	saleID := vars["order_id"]
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	}

	if err := database.DB.Delete(&sale).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}
// Method: POST
func AddNewGateInOrder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	saleID := vars["order_id"]
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	}

	var gateRowNumber int64

	err = r.ParseMultipartForm(256)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}
	gateType := models.GateType(r.FormValue("gateType"))
	switch gateType {
	case models.GateTypeInd:
		cfg, err := getGateCfg(gateType, true)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		gate := models.SalesAndGate{
			SaleID:        sale.ID,
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

		if err := database.DB.Clauses(clause.Returning{}).Create(&gate).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		if err := database.DB.Create(&models.IndustrialGatesAndSalesDrive{
			SaleID:    sale.ID,
			RowNumber: gate.RowNumber,
			DriveID:   int32(cfg.IndustrialDrives[0].ID),
		}).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		gateRowNumber = gate.RowNumber

	case models.GateTypeRes:
		cfg, err := getGateCfg(gateType, true)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		gate := models.SalesAndGate{
			SaleID:        sale.ID,
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

		if err := database.DB.Clauses(clause.Returning{}).Create(&gate).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		if err := database.DB.Create(&models.ResidentialGatesAndSalesDriveRail{
			SaleID:    sale.ID,
			RowNumber: gate.RowNumber,
			DriveID:   int32(cfg.ResidentialDrives[0].ID),
			RailID:    int32(cfg.Rails[0].ID),
		}).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		gateRowNumber = gate.RowNumber

	default:
		helpers.WriteError(w, errors.New("Неправильный тип ворот"), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/%d", gateRowNumber))
	w.WriteHeader(http.StatusCreated)
}

// Route: /orders
// Method: POST
func CreateNewOrder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	role := user.Role.Name

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var orderData proto_files.OrderRequest
	if err := proto.Unmarshal(data, &orderData); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	// Проверка на валидность типов ворот в заказе
	for _, gate := range orderData.OrderGates {
		_, err := models.GateTypeFromProto(gate.GateType)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
	}

	var order models.Sale
	if role == "dealer" {
		var managerId int64
		if err := database.DB.Model(&models.ManagerAndDealer{}).
			Select("manager_id").
			Where("dealer_id = ?", user.ID).
			Find(&managerId).Error; err != nil {
			helpers.WriteError(w, err, http.StatusNotFound)
			return
		}

		order = models.Sale{
			ClientID:  &user.ID,
			ManagerID: managerId,
			Status:    models.OrderStatusNew,
		}

		if err := database.DB.Create(&order).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		order = models.Sale{
			ClientID:  nil,
			ManagerID: user.ID,
			Status:    models.OrderStatusNew,
		}

		if err := database.DB.Create(&order).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}

	for productId, amount := range orderData.Products {
		salesAndProduct := models.SalesAndProduct{
			SaleID:    order.ID,
			ProductID: productId,
			Amount:    amount,
		}

		if err := database.DB.Create(&salesAndProduct).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}

	var wg sync.WaitGroup
	for i, gate := range orderData.OrderGates {
		gateType, _ := models.GateTypeFromProto(gate.GateType)
		driveType := models.GetDriveTypeFromProto(gate.Drive)
		if driveType == "unknown" {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		orderDetails := models.SalesAndGate{
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

		if err := database.DB.Create(&orderDetails).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		switch d := gate.Drive.DriveType.(type) {
		case *proto_files.Drive_Industrial:
			driveID := d.Industrial.DriveId
			indGatesAndSalesDrive := models.IndustrialGatesAndSalesDrive{
				SaleID:    order.ID,
				RowNumber: int64(i + 1),
				DriveID:   int32(driveID),
			}

			if err := database.DB.Create(&indGatesAndSalesDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		case *proto_files.Drive_Residential:
			driveID := d.Residential.DriveId
			railID := d.Residential.RailId
			resGatesAndSalesDriveRail := models.ResidentialGatesAndSalesDriveRail{
				SaleID:    order.ID,
				RowNumber: int64(i + 1),
				DriveID:   int32(driveID),
				RailID:    int32(railID),
			}

			if err := database.DB.Create(&resGatesAndSalesDriveRail).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		case *proto_files.Drive_Manual:
			chain := d.Manual.ChainLength
			gatesAndSalesManualDrive := models.GatesAndSalesManualDrive{
				SaleID:      order.ID,
				RowNumber:   int64(i + 1),
				ChainLength: chain,
			}
			if err := database.DB.Create(&gatesAndSalesManualDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		default:
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		if len(gate.Options) == 0 {
			continue
		}

		var gateAndSalesOptions []models.GatesAndSalesOption
		for optionId, amount := range gate.Options {
			gateAndSalesOption := models.GatesAndSalesOption{
				SaleID:    order.ID,
				RowNumber: int64(i + 1),
				OptionID:  optionId,
				Amount:    amount,
			}

			gateAndSalesOptions = append(gateAndSalesOptions, gateAndSalesOption)
		}

		if err := database.DB.Create(&gateAndSalesOptions).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}

	wg.Wait()
	w.WriteHeader(http.StatusCreated)
}

// Route: /orders/{order_id}/products
// Method: POST
func AddNewProductInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/products/{product_id}
// Method: PUT
func UpdateProductList(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/products/{product_id}
// Method: DELETE
func DeleteProductFromOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/{gate_id}
// Method: DELETE
func DeleteGateFromOrder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	saleID := vars["order_id"]
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	}
	rowNumber := vars["gate_id"]

	var gate models.SalesAndGate
	userRole := user.Role.Name
	switch userRole {
	case "dealer":
		if err := database.DB.
			Table("sales_and_gates AS sag").
			Joins("JOIN sales AS s ON s.id = sag.sale_id").
			Where("sag.sale_id = ? AND sag.row_number = ? AND s.client_id = ?", sale.ID, rowNumber, user.ID).
			First(&gate).Error; err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
	case "manager":
		if err := database.DB.
			Table("sales_and_gates AS sag").
			Joins("JOIN sales AS s ON s.id = sag.sale_id").
			Where("sag.sale_id = ? AND sag.row_number = ? AND s.manager_id = ?", sale.ID, rowNumber, user.ID).
			First(&gate).Error; err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

	default:
		helpers.WriteError(w, errors.New("Данная роль не имеет доступа к удалению записи"), http.StatusForbidden)
		return
	}

	if err := database.DB.Delete(&gate).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: PUT
func UpdateGateInOrder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	saleID := vars["order_id"]
	sale, err := getAccessibleSale(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	}

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var gateData proto_files.GateConfig
	if err := proto.Unmarshal(data, &gateData); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateId, err := strconv.Atoi(vars["gate_id"])
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var gate models.SalesAndGate
	if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&gate).Error; err != nil {
		helpers.WriteError(w, err, http.StatusNotFound)
		return
	}

	gate.Width = gateData.Width
	gate.Height = gateData.Height
	gate.Headroom = gateData.Headroom
	gate.LiftTypeID = gateData.LiftTypeId
	gate.ColorOutID = gateData.ColorOutId
	gate.CycleAmountID = gateData.CycleAmountId
	gate.Amount = gateData.Amount

	switch gateData.Drive.DriveType.(type) {
	case *proto_files.Drive_Industrial:
		if gate.DriveType != models.IndDriveType {
			indGatesAndSalesDrive := models.IndustrialGatesAndSalesDrive{
				SaleID:    int64(sale.ID),
				RowNumber: int64(gateId),
				DriveID:   int32(gateData.Drive.GetIndustrial().DriveId),
			}
			if err := database.DB.Create(&indGatesAndSalesDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}

			switch gate.DriveType {
			case models.ResDriveType:
				var resGatesAndSalesDriveRail models.ResidentialGatesAndSalesDriveRail
				if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&resGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusNotFound)
					return
				}

				if err := database.DB.Delete(&resGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusInternalServerError)
					return
				}
			case models.ManualDriveType:
				var gatesAndSalesManualDrive models.GatesAndSalesManualDrive
				if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&gatesAndSalesManualDrive).Error; err != nil {
					helpers.WriteError(w, err, http.StatusNotFound)
					return
				}

				if err := database.DB.Delete(&gatesAndSalesManualDrive).Error; err != nil {
					helpers.WriteError(w, err, http.StatusInternalServerError)
					return
				}
			}
		} else {
			var indGatesAndSalesDrive models.IndustrialGatesAndSalesDrive
			if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&indGatesAndSalesDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusNotFound)
				return
			}

			indGatesAndSalesDrive.DriveID = int32(gateData.Drive.GetIndustrial().DriveId)

			if err := database.DB.Save(&indGatesAndSalesDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

	case *proto_files.Drive_Residential:
		if gate.DriveType != models.ResDriveType {
			resGatesAndSalesDriveRail := models.ResidentialGatesAndSalesDriveRail{
				SaleID:    int64(sale.ID),
				RowNumber: int64(gateId),
				DriveID:   int32(gateData.Drive.GetResidential().DriveId),
				RailID:    int32(gateData.Drive.GetResidential().RailId),
			}
			if err := database.DB.Create(&resGatesAndSalesDriveRail).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}

			switch gate.DriveType {
			case models.IndDriveType:
				var indGatesAndSalesDriveRail models.IndustrialGatesAndSalesDrive
				if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&indGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusNotFound)
					return
				}
				if err := database.DB.Delete(&indGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusInternalServerError)
					return
				}

			case models.ManualDriveType:
				var gatesAndSalesManualDrive models.GatesAndSalesManualDrive
				if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&gatesAndSalesManualDrive).Error; err != nil {
					helpers.WriteError(w, err, http.StatusNotFound)
					return
				}
				if err := database.DB.Delete(&gatesAndSalesManualDrive).Error; err != nil {
					helpers.WriteError(w, err, http.StatusInternalServerError)
					return
				}
			}
		} else {
			var resGatesAndSalesDriveRail models.ResidentialGatesAndSalesDriveRail
			if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&resGatesAndSalesDriveRail).Error; err != nil {
				helpers.WriteError(w, err, http.StatusNotFound)
				return
			}
			resGatesAndSalesDriveRail.DriveID = int32(gateData.Drive.GetResidential().DriveId)
			resGatesAndSalesDriveRail.RailID = int32(gateData.Drive.GetResidential().RailId)
			if err := database.DB.Save(&resGatesAndSalesDriveRail).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}
	case *proto_files.Drive_Manual:
		if gate.DriveType != models.ManualDriveType {
			gatesAndSalesManualDrive := models.GatesAndSalesManualDrive{
				SaleID:      int64(sale.ID),
				RowNumber:   int64(gateId),
				ChainLength: gateData.Drive.GetManual().ChainLength,
			}
			if err := database.DB.Create(&gatesAndSalesManualDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}

			switch gate.DriveType {
			case models.ResDriveType:
				var resGatesAndSalesDriveRail models.ResidentialGatesAndSalesDriveRail
				if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&resGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusNotFound)
					return
				}
				if err := database.DB.Delete(&resGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusInternalServerError)
					return
				}
			case models.IndDriveType:
				var indGatesAndSalesDriveRail models.IndustrialGatesAndSalesDrive
				if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&indGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusNotFound)
					return
				}
				if err := database.DB.Delete(&indGatesAndSalesDriveRail).Error; err != nil {
					helpers.WriteError(w, err, http.StatusInternalServerError)
					return
				}
			}
		} else {
			var gatesAndSalesManualDrive models.GatesAndSalesManualDrive
			if err := database.DB.Where("sale_id = ? AND row_number = ?", sale.ID, gateId).First(&gatesAndSalesManualDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusNotFound)
				return
			}
			gatesAndSalesManualDrive.ChainLength = gateData.Drive.GetManual().ChainLength
			if err := database.DB.Save(&gatesAndSalesManualDrive).Error; err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}
	default:
		helpers.WriteError(w, errors.New("Неизвестный тип привода"), http.StatusBadRequest)
		return
	}

	gate.DriveType = models.GetDriveTypeFromProto(gateData.Drive)

	if err := database.DB.Save(&gate).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
