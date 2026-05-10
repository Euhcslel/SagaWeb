package http

import (
	"errors"
	"net/http"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"strconv"
)

// Route: /user
// Method: GET
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	userInfo, err := service.GetUserInfo(user)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":      "user.css",
		"userInfo": userInfo,
		"user":     user,
	}

	if err := templates.ExecuteTemplate(w, "info.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /user/dealers
// Method: GET
func GetUserDealers(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	dealers, err := service.GetUserDealers(user)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":     "",
		"user":    user,
		"dealers": dealers,
	}

	if err := templates.ExecuteTemplate(w, "dealers.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /user
// Method: POST
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userInfo := types.UpdatedUserInfo{
		Fullname: r.FormValue("fullname"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
	}

	if user.Role == enums.DealerRole {
		userInfo.Company = r.FormValue("company")
		userInfo.Address = r.FormValue("address")
	}

	if err := service.UpdateUserInfo(user, userInfo); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

// Route: /dealers/requests
// Method: GET
func GetDealersRegRequests(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	regRequests, err := service.GetDealersRegistrationRequests()
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"user":        user,
		"css":         "",
		"regRequests": regRequests,
	}
	if err := templates.ExecuteTemplate(w, "reg_requests.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /dealers/requests/{request_id}/confirm
// Method: POST
func ConfirmDealerRegRequest(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	requestId, err := strconv.ParseInt(r.PathValue("request_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := service.ConfirmDealerRegistrationRequest(user, requestId); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /dealers/requests/{request_id}/reject
// Method: POST
func RejectDealerRegRequest(w http.ResponseWriter, r *http.Request) {
	requestId, err := strconv.ParseInt(r.PathValue("request_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := service.RejectDealerRegistrationRequest(requestId); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
