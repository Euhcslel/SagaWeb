package http

import (
	"errors"
	"net/http"
	errs "project/internal/errors"
	"project/internal/helpers"
	"project/internal/service"
	"project/internal/utils"
)

// Route: /user
// Method: GET
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	data := map[string]any{
		"css":  "user.css",
		"user": user,
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
