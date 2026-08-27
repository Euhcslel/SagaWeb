package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/Euhcslel/SagaWeb/internal/utils"
)

// GetUserInfo возвращает страницу с информацией о текущем пользователе.
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

// GetUserDealers возвращает странницу со списком дилеров пользователя.
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
		"css":     "dealers.css",
		"user":    user,
		"dealers": dealers,
	}

	if err := templates.ExecuteTemplate(w, "dealers.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// UpdateUserInfo обновляет информацию о пользователе.
// Route: /user
// Method: POST
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	userInfo := types.UpdatedUserInfo{
		Fullname:    r.FormValue("fullname"),
		Email:       r.FormValue("email"),
		Phone:       r.FormValue("phone"),
		NewPassword: r.FormValue("new_password"),
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

// GetDealersRegRequests возвращает страницу с заявками на регистрацию от дилеров.
// Route: /dealers/requests
// Method: GET
func GetDealersRegRequests(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if user.Role != enums.ManagerRole && user.Role != enums.AdminRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

	regRequests, err := service.GetDealersRegistrationRequests()
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"user":        user,
		"css":         "user.css",
		"regRequests": regRequests,
	}
	if err := templates.ExecuteTemplate(w, "reg_requests.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// ConfirmDealerRegRequest подтверждает указаннную заявку на регистрацию от дилера.
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
		if errors.Is(err, errs.ErrForbidden) {
			helpers.WriteError(w, err, http.StatusForbidden)
			return
		}
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// RejectDealerRegRequest отклоняет указаннную заявку на регистрацию от дилера.
// Route: /dealers/requests/{request_id}/reject
// Method: POST
func RejectDealerRegRequest(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if user.Role != enums.ManagerRole && user.Role != enums.AdminRole {
		helpers.WriteError(w, errs.ErrForbidden, http.StatusForbidden)
		return
	}

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

// UpdateCompanyDetails обновляет информацию о компании пользователя.
// Route: /user/company
// Method: POST
func UpdateCompanyDetails(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err := service.UpdateCompanyDetails(user, r.FormValue("inn"), r.FormValue("kpp"), r.FormValue("ogrn"))
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

// GenerateDealerContract подставляет в шаблон данные о дилере и возвращает сгенерированный договор.
// Route: /user/dealers/{dealer_id}/contract/generate
// Method: GET
func GenerateDealerContract(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	dealerID, err := strconv.ParseInt(r.PathValue("dealer_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	data, contractNumber, err := service.GenerateDealerContract(user, dealerID)
	if errors.Is(err, errs.ErrForbidden) {
		helpers.WriteError(w, err, http.StatusForbidden)
		return
	} else if errors.Is(err, errs.ErrIncompleteCompanyDetails) {
		helpers.WriteError(w, err, http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="contract_%d.docx"`, contractNumber))
	w.Write(data)
}

// AttachContractToDealer прикрепляет файл договора к указанному дилеру.
// Route: /user/dealers/{dealer_id}/contract
// Method: POST
func AttachContractToDealer(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	dealerID, err := strconv.ParseInt(r.PathValue("dealer_id"), 10, 64)
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	signedAt, err := time.Parse("2006-01-02", r.FormValue("signed_at"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("contract_file")
	if err != nil {
		helpers.WriteError(w, fmt.Errorf("файл договора обязателен"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := service.AttachContractToDealer(user, dealerID, r.FormValue("contract_number"), signedAt, file, handler); err != nil {
		if errors.Is(err, errs.ErrForbidden) {
			helpers.WriteError(w, err, http.StatusForbidden)
			return
		}
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/dealers", http.StatusSeeOther)
}
