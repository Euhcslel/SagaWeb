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
	"github.com/Euhcslel/SagaWeb/internal/domain/manual_drive_prices"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	rails_domain "github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

// Route: /tables/{table_name}
// Method: GET
func GetDataBaseRedactor(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	tablePageData, err := service.GetTablePageData(tableName)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidTableType) {
			helpers.WriteError(w, errs.ErrNotFound, http.StatusNotFound)
			return
		}
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":       "admin.css",
		"user":      user,
		"tableData": tablePageData,
	}

	if tableName == "users" {
		var roles []enums.Role
		for _, r := range enums.GetAllRoles() {
			if r != enums.ClientRole {
				roles = append(roles, r)
			}
		}
		data["roles"] = roles
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
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	tablesWithImagesUpdate := map[string]bool{
		"products": true, "options": true,
		"industrial_drives": true, "residential_drives": true, "rails": true,
	}
	if tablesWithImagesUpdate[tableName] {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
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

		imagePath := r.FormValue("image-path")
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = industrial_gate_drives.IndustrialGateDrive{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
			ImagePath:      imagePath,
		}
	case "lift_types":
		maxHeadroom, err := strconv.ParseInt(r.FormValue("max-headroom"), 10, 32)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		minHeadroom, err := strconv.ParseInt(r.FormValue("min-headroom"), 10, 32)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
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
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		retailPrice, err := decimal.NewFromString(r.FormValue("retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}

		imagePath := r.FormValue("image-path")
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = options.Option{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			ImagePath:      imagePath,
		}
	case "products":
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

		imagePath := r.FormValue("image-path")
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = products.Product{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			ImagePath:      imagePath,
		}
	case "rails":
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

		imagePath := r.FormValue("image-path")
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = rails_domain.Rail{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
			ImagePath:      imagePath,
		}
	case "residential_drives":
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

		imagePath := r.FormValue("image-path")
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = residential_gate_drives.ResidentialGateDrive{
			ID:             rowId,
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
			ImagePath:      imagePath,
		}
	case "users":
		role := enums.Role(r.FormValue("role"))
		valid := false
		for _, r := range enums.GetAllRoles() {
			if r == role {
				valid = true
				break
			}
		}

		if !valid {
			helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
			return
		}

		tableData = users.User{
			ID:          rowId,
			Fullname:    r.FormValue("fullname"),
			Email:       r.FormValue("email"),
			PhoneNumber: r.FormValue("phone"),
			Role:        role,
		}

		if pw := r.FormValue("password"); pw != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
			if err := repository.AdminUpdateUserPassword(rowId, hash); err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}
	case "manual_drive_prices":
		chainMeterRetailPrice, err := decimal.NewFromString(r.FormValue("chain-meter-retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		chainMeterWholesalePrice, err := decimal.NewFromString(r.FormValue("chain-meter-wholesale-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		rcpRetailPrice, err := decimal.NewFromString(r.FormValue("rcp-retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		rcpWholesalePrice, err := decimal.NewFromString(r.FormValue("rcp-wholesale-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		tableData = manual_drive_prices.ManualDrivePrice{
			ID:                       rowId,
			ChainMeterRetailPrice:    chainMeterRetailPrice,
			ChainMeterWholesalePrice: chainMeterWholesalePrice,
			RcpRetailPrice:           rcpRetailPrice,
			RcpWholesalePrice:        rcpWholesalePrice,
		}
	default:
		helpers.WriteError(w, errs.ErrInvalidTableType, http.StatusBadRequest)
		return
	}

	if err := service.UpdateRow(tableName, tableData); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Route: /tables/{table_name}/{row_id}
// Method: DELETE
func DeleteRow(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	role := user.Role
	if role != enums.AdminRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
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
		if errors.Is(err, errs.ErrInvalidTableType) {
			helpers.WriteError(w, errs.ErrNotFound, http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrForeignKeyViolation) {
			helpers.WriteError(w, err, http.StatusConflict)
			return
		}
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
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
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
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	tableName := r.PathValue("table_name")
	if tableName == "" {
		helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
		return
	}

	tablesWithImages := map[string]bool{
		"products": true, "options": true,
		"industrial_drives": true, "residential_drives": true, "rails": true,
	}
	if tablesWithImages[tableName] {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
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

		imagePath := "no_image"
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = industrial_gate_drives.IndustrialGateDrive{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
			ImagePath:      imagePath,
		}
	case "lift_types":
		maxHeadroom, err := strconv.ParseInt(r.FormValue("max-headroom"), 10, 32)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		minHeadroom, err := strconv.ParseInt(r.FormValue("min-headroom"), 10, 32)
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
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

		tableData = lift_types.LiftType{
			Name:            r.FormValue("name"),
			MinHeadroom:     int32(minHeadroom),
			MaxHeadroom:     int32(maxHeadroom),
			WholesaleMarkup: wholesaleMakup,
			RetailMarkup:    retailMakup,
		}
	case "options":
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

		imagePath := "no_image"
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = options.Option{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			ImagePath:      imagePath,
		}
	case "products":
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

		imagePath := "no_image"
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = products.Product{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			ImagePath:      imagePath,
		}
	case "residential_drives":
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

		imagePath := "no_image"
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = residential_gate_drives.ResidentialGateDrive{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
			ImagePath:      imagePath,
		}
	case "rails":
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

		imagePath := "no_image"
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err = service.UploadProductImage(file, handler)
			if err != nil {
				helpers.WriteError(w, err, http.StatusInternalServerError)
				return
			}
		}

		tableData = rails_domain.Rail{
			Name:           r.FormValue("name"),
			WholesalePrice: wholesalePrice,
			RetailPrice:    retailPrice,
			Specifications: specifications,
			ImagePath:      imagePath,
		}
	case "users":
		role := enums.Role(r.FormValue("role"))
		valid := false
		for _, r := range enums.GetAllRoles() {
			if r == role {
				valid = true
				break
			}
		}
		if !valid {
			helpers.WriteError(w, errs.ErrBadRequest, http.StatusBadRequest)
			return
		}

		passwordRaw := r.FormValue("password")
		if passwordRaw == "" {
			helpers.WriteError(w, errors.New("пароль обязателен"), http.StatusBadRequest)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
		if err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		tableData = users.User{
			Fullname:     r.FormValue("fullname"),
			Email:        r.FormValue("email"),
			PhoneNumber:  r.FormValue("phone"),
			Role:         role,
			PasswordHash: hash,
		}
	case "manual_drive_prices":
		chainMeterRetailPrice, err := decimal.NewFromString(r.FormValue("chain-meter-retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		chainMeterWholesalePrice, err := decimal.NewFromString(r.FormValue("chain-meter-wholesale-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		rcpRetailPrice, err := decimal.NewFromString(r.FormValue("rcp-retail-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		rcpWholesalePrice, err := decimal.NewFromString(r.FormValue("rcp-wholesale-price"))
		if err != nil {
			helpers.WriteError(w, err, http.StatusBadRequest)
			return
		}
		tableData = manual_drive_prices.ManualDrivePrice{
			ChainMeterRetailPrice:    chainMeterRetailPrice,
			ChainMeterWholesalePrice: chainMeterWholesalePrice,
			RcpRetailPrice:           rcpRetailPrice,
			RcpWholesalePrice:        rcpWholesalePrice,
		}
	default:
		helpers.WriteError(w, errs.ErrInvalidTableType, http.StatusBadRequest)
		return
	}

	err := service.AddNewDataBaseTableRow(tableName, tableData)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tables/"+tableName, http.StatusSeeOther)
}
