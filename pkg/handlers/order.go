package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"project/pkg/types"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Route: /orders
// Method: GET
func GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)

	var orders []models.Sale
	err = database.DB.
		Preload("Client").
		Preload("Manager").
		Where("client_id = ?", user.ID).
		Limit(5).
		Find(&orders).Error
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":    "",
		"orders": orders,
	}

	if err := templates.ExecuteTemplate(w, "list.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
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

	// Получаем все ворота выбранного заказа
	var orderGates []models.SalesAndGate
	database.DB.Preload("GateType").
		Preload("LiftType").
		Preload("ColorIn").
		Preload("ColorOut").
		Preload("CycleAmount").
		Preload("Status").
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
		Products: []models.Product{},
	}
	for _, og := range orderGates {
		order.Gates = append(order.Gates, types.Gate{
			Gate:    og,
			Options: optionsMap[og.RowNumber],
		})
	}

	data := map[string]any{
		"css":    "",
		"orders": order,
	}

	if err := templates.ExecuteTemplate(w, "order.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: GET
func GetGateInOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var gate models.SalesAndGate
	if err := database.DB.Model(models.SalesAndGate{}).
	Preload("GateType").
	Preload("LiftType").
	Preload("ColorIn").
	Preload("ColorOut").
	Preload("CycleAmount").
	Preload("Status").
	Where("sale_id = ? and row_number = ?", vars["order_id"], vars["gate_id"]).
	Find(&gate).Error; err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	var options []models.GatesAndSalesOption
	if err := database.DB.Model(models.GatesAndSalesOption{}).Preload("Option").Where("sale_id = ? and row_number = ?", vars["order_id"], vars["gate_id"]).Find(&options).Error; err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	stringGateType := strconv.Itoa(int(gate.GateTypeID))
	cfg, err := getGateCfg(stringGateType)
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":     "",
		"gate":     gate,
		"options": options,
		"cfg": cfg,
	}

	if err := templates.ExecuteTemplate(w, "gate.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
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
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)
	role := user.Role.Name

	err = r.ParseMultipartForm(64 << 20)
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusBadRequest)
		return
	}

	var order models.Sale
	if role == "client" || role == "dealer" {
		var managerID int64
		database.DB.Model(&models.ManagerAndDealer{}).
			Select("manager_id").
			Where("dealer_id = ?", user.ID).
			Find(&managerID)

		order = models.Sale{
			ClientID: &user.ID,
		}

		if managerID > 0 {
			order.ManagerID = &managerID
		}
		database.DB.Create(&order)
	}


	orderGatesJSON := r.FormValue("orderGates")
	productsJSON := r.FormValue("products")

	var orderRequest types.OrderRequest

	if err := json.Unmarshal([]byte(orderGatesJSON), &orderRequest.OrderGates); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse orderGates: %v", err), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(productsJSON), &orderRequest.Products); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse products: %v", err), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	for i, gate := range orderRequest.OrderGates {
		gateTypeInt, _ := strconv.ParseInt(gate.GateType, 10, 32)
		width := gate.Width
		height := gate.Height
		liftTypeID, _ := strconv.ParseInt(gate.LiftType, 10, 64)
		colorInID, _ := strconv.ParseInt(gate.ColorIn, 10, 64)
		colorOutID, _ := strconv.ParseInt(gate.ColorOut, 10, 64)
		cycleAmountID, _ := strconv.ParseInt(gate.CycleAmount, 10, 64)

		orderDetails := models.SalesAndGate{
			SaleID:        order.ID,
			RowNumber:     int64(i + 1),
			GateTypeID:    int32(gateTypeInt),
			Width:         int32(width),
			Height:        int32(height),
			LiftTypeID:    liftTypeID,
			ColorInID:     colorInID,
			ColorOutID:    colorOutID,
			CycleAmountID: cycleAmountID,
			TotalPrice:    int32(gate.GatePrice),
			StatusID:      9,
		}

		wg.Go(func() {
			database.DB.Create(&orderDetails)
		})

		if len(gate.Options) == 0 {continue}
		var gateAndSalesOptions []models.GatesAndSalesOption
		for _, option := range gate.Options {
			optionId, _ := strconv.ParseInt(option, 10, 32)
			gateAndSalesOption := models.GatesAndSalesOption{
				SaleID:    order.ID,
				RowNumber: int64(i + 1),
				OptionID:  int32(optionId),
				Amount:    1,
			}

			gateAndSalesOptions = append(gateAndSalesOptions, gateAndSalesOption)
		}

		wg.Go(func() {
			database.DB.Create(&gateAndSalesOptions)
		})
	}

	wg.Wait()
	http.Redirect(w, r, "/orders", http.StatusSeeOther)
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
