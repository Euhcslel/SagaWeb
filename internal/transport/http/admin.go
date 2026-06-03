package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"github.com/shopspring/decimal"
)

// Route: /tables/{table_name}
// Method: GET
func GetDataBaseRedactor(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	tablePageData, err := service.GetTablePageData(tableName)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":       "admin.css",
		"user":      user,
		"tableData": tablePageData,
	}

	if err := templates.ExecuteTemplate(w, tableName+".html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /tables/{table_name}
// Method: PUT
func UpdateRow(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var tableData any
	rowId, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}
	switch tableName {
	case "colors":
		tableData = colors.Color{
			ID:   rowId,
			Code: r.FormValue("code"),
			Hex:  r.FormValue("hex"),
		}
	case "cycle_amounts":
		wholesaleMakup, err := decimal.NewFromString(r.FormValue("wholesale-markup"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		retailMakup, err := decimal.NewFromString(r.FormValue("retail-markup"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = cycle_amounts.CycleAmount{
			ID:              rowId,
			Amount:          r.FormValue("amount"),
			WholesaleMarkup: wholesaleMakup,
			RetailMarkup:    retailMakup,
		}
	case "industrial_drives":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		specificationsValue := r.FormValue("specifications")
		specifications := &specificationsValue

		tableData = industrial_gate_drives.IndustrialGateDrive{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
		}
	case "lift_types":
		maxHeadroom, err := strconv.ParseInt(r.FormValue("max-headroom"), 10, 32)
		minHeadroom, err := strconv.ParseInt(r.FormValue("min-headroom"), 10, 32)
		wholesaleMakup, err := decimal.NewFromString(r.FormValue("wholesale-markup"))
		retailMakup, err := decimal.NewFromString(r.FormValue("retail-markup"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = lift_types.LiftType{
			ID:              rowId,
			Name:            r.FormValue("name"),
			MinHeadroom:     int32(minHeadroom),
			MaxHeadroom:     int32(maxHeadroom),
			WholesaleMarkup: wholesaleMakup,
			RetailMarkup:    retailMakup,
		}
	case "options":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = options.Option{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
		}
	case "products":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = products.Product{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
		}
	case "residential_drives":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		specificationsValue := r.FormValue("specifications")
		specifications := &specificationsValue

		tableData = residential_gate_drives.ResidentialGateDrive{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
		}
	default:
		helpers.WriteError(w, errors.New("invalid table name"), http.StatusBadRequest)
		return
	}

	if err := service.UpdateRow(tableName, tableData); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /tables/{table_name}/{row_id}
// Method: DELETE
func DeleteRow(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	rowId, err := strconv.ParseInt(r.PathValue("row_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := service.DeleteRow(tableName, rowId); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /tables
// Method: GET
func GetDataBaseTableList(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	data := map[string]any{
		"css":  "admin.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "tables.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /tables/{table_name}
// Method: POST
func AddNewDataBaseTableRow(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var tableData any
	switch tableName {
	case "colors":
		tableData = colors.Color{
			Code: r.FormValue("code"),
			Hex:  r.FormValue("hex"),
		}
	case "cycle_amounts":
		wholesaleMakup, err := decimal.NewFromString(r.FormValue("wholesale-markup"))
		retailMakup, err := decimal.NewFromString(r.FormValue("retail-markup"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = cycle_amounts.CycleAmount{
			Amount:          r.FormValue("amount"),
			WholesaleMarkup: wholesaleMakup,
			RetailMarkup:    retailMakup,
		}
	case "industrial_drives":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		specificationsValue := r.FormValue("specifications")
		specifications := &specificationsValue

		tableData = industrial_gate_drives.IndustrialGateDrive{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
		}
	case "lift_types":
		maxHeadroom, err := strconv.ParseInt(r.FormValue("max-headroom"), 10, 32)
		minHeadroom, err := strconv.ParseInt(r.FormValue("min-headroom"), 10, 32)
		wholesaleMakup, err := decimal.NewFromString(r.FormValue("wholesale-markup"))
		retailMakup, err := decimal.NewFromString(r.FormValue("retail-markup"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = lift_types.LiftType{
			Name:            r.FormValue("name"),
			MinHeadroom:     int32(minHeadroom),
			MaxHeadroom:     int32(maxHeadroom),
			WholesaleMarkup: wholesaleMakup,
			RetailMarkup:    retailMakup,
		}
	case "options":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = options.Option{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
		}
	case "products":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		tableData = products.Product{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
		}
	case "residential_drives":
		wholesalePrice, err := decimal.NewFromString(r.FormValue("wholesale-price"))
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		specificationsValue := r.FormValue("specifications")
		specifications := &specificationsValue

		tableData = residential_gate_drives.ResidentialGateDrive{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
		}
	default:
		helpers.WriteError(w, errors.New("invalid table name"), http.StatusBadRequest)
		return
	}

	err := service.AddNewDataBaseTableRow(tableName, tableData)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tables/"+tableName, http.StatusSeeOther)
}
