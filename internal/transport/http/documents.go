package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Euhcslel/SagaWeb/internal/generated"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"google.golang.org/protobuf/proto"
)

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

// Route: /offer
// Method: POST
func GetOfferForOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var orderData generated.OrderRequest

	err = proto.Unmarshal(body, &orderData)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	file, err := service.GetOfferForOrder(
		user,
		&orderData,
	)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="offer.pdf"`)
	_, err = w.Write(file)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
