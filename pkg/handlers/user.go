package handlers

import (
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
)

// Route: /user
// Method: GET
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	} else {
		token := sessionToken.Value
		user, err = helpers.GetUserBySessionToken(token)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
	}

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
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user, err := helpers.GetUserBySessionToken(token)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	role := user.Role.Name

	if role == "manager" {
		var dealers []models.ManagerAndDealer
		if err := database.DB.Model(models.ManagerAndDealer{}).Preload("Dealer").Where("manager_id = ?", user.ID).Find(&dealers).Error; err != nil {
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
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
	}

}
