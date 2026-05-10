package http

import (
	"errors"
	"fmt"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/generated"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"io"
	"net/http"
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
	//To do
}

// Route: /orders/{order_id}
// Method: GET
func GetUserOrderById(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	pageData, err := service.GetOrderPageData(user, orderID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	statuses := enums.GetAllOrderStatuses()

	data := map[string]any{
		"css":           "/order.css",
		"products":      pageData.Products,
		"order":         pageData.Order,
		"user":          user,
		"orderStatuses": statuses,
		"orderId":       orderID,
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

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateID, err := strconv.ParseInt(r.PathValue("gate_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	pageData, err := service.GetCurrentGatePageData(user, orderID, gateID)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":              "gate.css",
		"gate":             pageData.Gate,
		"options":          pageData.Options,
		"cfg":              pageData.Configuration,
		"industrialDrive":  pageData.IndustrialDrive,
		"residentialDrive": pageData.ResidentialDrive,
		"manualDrive":      pageData.ManualDrive,
		"orderStatus":      pageData.OrderStatus,
		"user":             user,
	}

	if err := templates.ExecuteTemplate(w, "gate.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}
// Method: DELETE
func DeleteUserOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = service.DeleteUserOrder(user, orderID)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}
// Method: POST
func AddNewGateInOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
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

	newGate, err := service.AddNewGateInOrder(user, orderID, gateType)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if errors.Is(err, errs.ErrInvalidGateType) {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	} else if err != nil {
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
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Route: /orders/{order_id}/products
// Method: PUT
func UpdateProductList(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
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

	var updateProductsRequest generated.UpdateProductsRequest
	if err := proto.Unmarshal(data, &updateProductsRequest); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := service.UpdateProductsInOrder(user, orderID, &updateProductsRequest); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/status
// Method: PUT
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
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

	var updateStatusRequest generated.UpdateOrderStatusRequest
	if err := proto.Unmarshal(data, &updateStatusRequest); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = service.UpdateOrderStatus(user, orderID, &updateStatusRequest)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if errors.Is(err, errs.ErrInvalidOrderStatus) {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: DELETE
func DeleteGateFromOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateID, err := strconv.ParseInt(r.PathValue("gate_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = service.DeleteGateFromOrder(user, orderID, gateID)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: PUT
func UpdateGateInOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
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

	err = service.UpdateGateInOrder(user, orderID, gateID, &gateData)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/documents
// Method: GET
func GetOrderDocuments(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := service.GetAllOrderDocuments(user, orderID)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Route: /orders/{order_id}/offer
// Method: POST
func UploadOfferToOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = service.UploadOfferToOrder(user, orderID, file, handler)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/orders/"+fmt.Sprint(orderID), http.StatusSeeOther)
}

// Route: /orders/{order_id}/contract
// Method: POST
func UploadContractToOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = service.UploadContractToOrder(user, orderID, file, handler)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/orders/"+fmt.Sprint(orderID), http.StatusSeeOther)
}

// Route: /orders/{order_id}/bill
// Method: POST
func UploadBillToOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = service.UploadBillToOrder(user, orderID, file, handler)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/orders/"+fmt.Sprint(orderID), http.StatusSeeOther)
}

// Route: /orders/{order_id}/documents/{document_type}/{document_name}
// Method: GET
func DownloadOrderDocument(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	documentType := r.PathValue("document_type")
	documentName := r.PathValue("document_name")

	fileInfo, err := service.GetFileInfo(user, orderID, documentType, documentName)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.FileName)
	http.ServeFile(w, r, fileInfo.FilePath)
}

// Route: /orders/{order_id}/documents/{document_type}/{document_name}
// Method: DELETE
func DeleteOrderDocument(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	documentType := r.PathValue("document_type")
	documentName := r.PathValue("document_name")

	if err := service.DeleteOrderDocument(user, orderID, documentType, documentName); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
