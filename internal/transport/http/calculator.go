package http

import (
	"net/http"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"strconv"

	"google.golang.org/protobuf/proto"
)

// GetCalculatorForUser возвращает страницу калькулятора ворот.
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

	w.Header().Set("Vary", "Cookie")
	if err := templates.ExecuteTemplate(w, "calc.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// GetPriceBasedOnSize возвращает цены на ворота по указанным в запросе параметрам.
// Также делит ответ на дилерскую и розничную цены. Если клиент - дилер, то возвращает обе,
// Если клиент - не дилер, тогда возвращает только розничную.
// Route: /sizes
// Method: GET
func GetPriceBasedOnSize(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	query := r.URL.Query()
	width, err := strconv.ParseInt(query.Get("width"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	height, err := strconv.ParseInt(query.Get("height"), 10, 64)
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
