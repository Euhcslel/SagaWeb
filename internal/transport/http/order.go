package http

import (
	"fmt"
	"io"
	"net/http"
	"project/internal/domain/enums"
	"project/internal/generated"
	"project/internal/helpers"
	"project/internal/service"
	"project/internal/utils"
	"strconv"

	"google.golang.org/protobuf/proto"
)

// Route: /orders
// Method: GET
func GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orders, err := service.GetAllUserOrders(user)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
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
	user := utils.UserFromContext(r.Context())

	saleID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	order, err := service.GetUserOrderByID(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":     "/order.css",
		"order":   order,
		"user":    user,
		"orderId": saleID,
	}

	if err := templates.ExecuteTemplate(w, "order.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: GET
func GetGateInOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	saleID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateID, err := strconv.ParseInt(r.PathValue("gate_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	pageData, err := service.GetCurrentGatePageData(user, saleID, gateID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":              "calc.css",
		"gate":             pageData.Gate,
		"options":          pageData.Options,
		"cfg":              pageData.Configuraion,
		"industrialDrive":  pageData.IndustrialDrive,
		"residentialDrive": pageData.ResidentialDrive,
		"manualDrive":      pageData.ManualDrive,
		"statuses":         enums.GetAllOrderStatuses(),
		"user":             user,
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
	user := utils.UserFromContext(r.Context())

	saleID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = service.DeleteUserOrder(user, saleID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}
// Method: POST
func AddNewGateInOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	saleID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(256)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}
	gateType := r.FormValue("gateType")

	newGate, err := service.AddNewGateInOrder(user, saleID, gateType)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/%d", newGate.RowNumber))
	w.WriteHeader(http.StatusCreated)
}

// Route: /orders
// Method: POST
func CreateNewOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var orderData generated.OrderRequest
	if err := proto.Unmarshal(data, &orderData); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = service.CreateNewOrder(user, &orderData)

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
	user := utils.UserFromContext(r.Context())

	saleID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateID, err := strconv.ParseInt(r.PathValue("gate_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := service.DeleteGateFromOrder(user, saleID, gateID); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: PUT
func UpdateGateInOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	saleID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateID, err := strconv.ParseInt(r.PathValue("gate_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var gateData generated.GateConfig
	if err := proto.Unmarshal(data, &gateData); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := service.UpdateGateInOrder(user, saleID, gateID, &gateData); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
