package handlers

import (
	"io"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"project/pkg/proto_files"
	"project/pkg/types"
	"sync"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

// Route: /orders
// Method: GET
func GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)
	if user == nil {
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
	vars := mux.Vars(r)
	orderId := vars["order_id"]

	user := helpers.GetUserBySessionToken(w, r)
	if user == nil {
		return
	}

	// Получаем все ворота выбранного заказа
	var orderGates []models.SalesAndGate
	database.DB.
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Model(models.SalesAndGate{}).
		Where("sale_id = ?", orderId).
		Find(&orderGates)

	// Собираем все ID ворот
	var gateIDs []int64
	for _, gate := range orderGates {
		gateIDs = append(gateIDs, gate.RowNumber)
	}

	// Получаем все опции всех ворот
	var options []models.GatesAndSalesOption
	database.DB.
		Where("row_number IN ?", gateIDs).
		Preload("Option").
		Find(&options)

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

	if err := database.DB.Preload("Product").Where("sale_id = ?", orderId).Find(&order.Products).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":   "/order.css",
		"order": order,
		"user":  user,
	}

	if err := templates.ExecuteTemplate(w, "order.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: GET
func GetGateInOrder(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)
	if user == nil {
		return
	}

	vars := mux.Vars(r)

	var gate models.SalesAndGate
	if err := database.DB.Model(models.SalesAndGate{}).
		Preload("LiftType").
		Preload("ColorOut").
		Preload("CycleAmount").
		Where("sale_id = ? and row_number = ?", vars["order_id"], vars["gate_id"]).
		Find(&gate).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var options []models.GatesAndSalesOption
	if err := database.DB.Model(models.GatesAndSalesOption{}).Preload("Option").Where("sale_id = ? and row_number = ?", vars["order_id"], vars["gate_id"]).Find(&options).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	cfg, err := getGateCfg(gate.GateType)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":      "gate.css",
		"gate":     gate,
		"options":  options,
		"cfg":      cfg,
		"statuses": models.GetAllOrderStatuses(),
		"user":     user,
	}

	if err := templates.ExecuteTemplate(w, "gate.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}/options
// Method: GET
func GetGateOptions(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/documents
// Method: GET
func GetOrderDocuments(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/products
// Method: GET
func GetProductsInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: DELETE
func DeleteUserOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: POST
func AddNewGateInOrder(w http.ResponseWriter, r *http.Request) {
	database.DB.Model(models.SalesAndGate{})
}

// Route: /orders
// Method: POST
func CreateNewOrder(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)
	if user == nil {
		return
	}
	role := user.Role.Name

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

		database.DB.Create(&order)
	} else {
		order = models.Sale{
			ClientID:  nil,
			ManagerID: user.ID,
			Status:    models.OrderStatusNew,
		}

		database.DB.Create(&order)
	}

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

	for productId, amount := range orderData.Products {
		salesAndProduct := models.SalesAndProduct{
			SaleID:    order.ID,
			ProductID: productId,
			Amount:    amount,
		}

		database.DB.Create(&salesAndProduct)
	}

	var wg sync.WaitGroup
	for i, gate := range orderData.OrderGates {
		gateType, _ := models.GateTypeFromProto(gate.GateType)

		orderDetails := models.SalesAndGate{
			SaleID:        order.ID,
			RowNumber:     int64(i + 1),
			GateType:      gateType,
			Width:         gate.Width,
			Height:        gate.Height,
			LiftTypeID:    gate.LiftTypeId,
			ColorOutID:    gate.ColorOutId,
			CycleAmountID: gate.CycleAmountId,
			TotalPrice:    gate.GatePrice,
		}

		wg.Go(func() {
			database.DB.Create(&orderDetails)
		})

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

		wg.Go(func() {
			database.DB.Create(&gateAndSalesOptions)
		})
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

}

// Route: /orders/{order_id}/{gate_id}
// Method: PUT
func UpdateGateInOrder(w http.ResponseWriter, r *http.Request) {

}
