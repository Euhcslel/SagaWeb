package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
)

// Route: /logistician/orders/{order_id}
// Method: PUT
func UpdateLogisticianOrder(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if user.Role != enums.LogisticianRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	orderID, err := strconv.ParseInt(r.PathValue("order_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	var body struct {
		ManufactureDate string `json:"manufacture_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	var manufactureDate *time.Time
	if body.ManufactureDate != "" {
		t, err := time.Parse("2006-01-02", body.ManufactureDate)
		if err != nil {
			helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
			return
		}
		manufactureDate = &t
	}

	if err := service.UpdateOrderByLogistician(user, orderID, manufactureDate); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Route: /logistician/size-import
// Method: GET
func GetSizeImportPage(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if user.Role != enums.LogisticianRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	data := map[string]any{
		"css":  "logistician.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "size_import.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /logistician/size-import
// Method: POST
func ImportSizes(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if user.Role != enums.LogisticianRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	indDealerFile, _, err := r.FormFile("ind_dealer")
	if err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer indDealerFile.Close()

	indClientFile, _, err := r.FormFile("ind_client")
	if err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer indClientFile.Close()

	resDealerFile, _, err := r.FormFile("res_dealer")
	if err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer resDealerFile.Close()

	resClientFile, _, err := r.FormFile("res_client")
	if err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer resClientFile.Close()

	if err := service.ImportSizePrices(indDealerFile, indClientFile, resDealerFile, resClientFile); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistician/size-import", http.StatusSeeOther)
}
