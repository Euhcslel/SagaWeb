package http

import (
	"net/http"
	"project/internal/domain/enums"
	"project/internal/helpers"
	"project/internal/service"
	"project/internal/utils"
	"strconv"

	"google.golang.org/protobuf/proto"
)

// Route: /calculator
// Method: GET
func GetCalculatorForUser(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	pageData, err := service.GetCalculatorPageData(user)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":                      "calc.css",
		"IndustrialConfiguration":  pageData.IndustrialConfiguration,
		"ResidentialConfiguration": pageData.ResidentialConfiguration,
		"user":                     user,
		"isDealer":                 pageData.IsDealer,
	}

	if err := templates.ExecuteTemplate(w, "calc.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /sizes
// Method: GET
func GetPriceBasedOnSize(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	query := r.URL.Query()
	width, err := strconv.Atoi(query.Get("width"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(query.Get("height"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateType, err := enums.DetermineGateType(query.Get("gateType"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := service.GetPriceBasedOnSize(width, height, gateType, user)
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
