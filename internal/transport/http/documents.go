package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	errs "github.com/Euhcslel/SagaWeb/internal/errors"
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
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
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
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
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
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
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
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
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
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if errors.Is(err, errs.ErrInvalidDocumentType) {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileInfo.FileName))
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

	err = service.DeleteOrderDocument(user, orderID, documentType, documentName)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if errors.Is(err, errs.ErrInvalidDocumentType) {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/appendice
// Method: GET
func GetAppendiceForOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	file, appendiceNumber, err := service.GetAppendiceForOrder(user, orderID)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if errors.Is(err, errs.ErrNoDealerContract) {
		helpers.WriteError(w, err, http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("appendice_%d_%d.pdf", orderID, appendiceNumber)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Write(file)
}

// Route: /orders/{order_id}/offer/client
// Method: GET
// Query: client_name (string)
func GetDealerOfferForClient(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	clientName := r.URL.Query().Get("client_name")
	if clientName == "" {
		helpers.WriteError(w, fmt.Errorf("client_name обязателен"), http.StatusBadRequest)
		return
	}

	file, err := service.GetDealerOfferForClient(user, orderID, clientName)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="offer_%d_client.pdf"`, orderID))
	w.Write(file)
}

// Route: /orders/{order_id}/offer/self
// Method: GET
// Query: client_name (string)
func GetDealerOfferForSelf(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	clientName := r.URL.Query().Get("client_name")
	if clientName == "" {
		helpers.WriteError(w, fmt.Errorf("client_name обязателен"), http.StatusBadRequest)
		return
	}

	file, err := service.GetDealerOfferForSelf(user, orderID, clientName)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="offer_%d_self.pdf"`, orderID))
	w.Write(file)
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
