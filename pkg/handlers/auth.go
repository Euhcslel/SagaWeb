package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Route: /log
// Method: GET
func GetSignInForm(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"css": cssPath + "auth.css",
	}

	if err := templates.ExecuteTemplate(w, "log.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /log
// Method: POST
func SignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user models.User
	err = database.DB.Model(&models.User{}).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}
	userId := user.ID
	dbPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	helpers.SetSession(w, userId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /reg
// Method: GET
func GetSignUpFrom(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"css": cssPath + "auth.css",
	}
	if err := templates.ExecuteTemplate(w, "reg.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /reg
// Method: POST
func SignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusBadRequest)
		return
	}

	accountType := r.FormValue("accountType")

	switch accountType {
	case "dealer":
		fullname := r.FormValue("fullname")
		company := r.FormValue("company")
		phone := r.FormValue("phone")
		email := r.FormValue("email")

		var status models.Status
		if err := database.DB.Model(&models.Status{}).Where("name = ?", "Ожидает подтверждения").Find(&status).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}

		request := models.DealerRegRequest{
			Company:     company,
			Fullname:    fullname,
			PhoneNumber: phone,
			Email:       email,
			StatusID:    int64(status.ID),
		}
		database.DB.Create(&request)

	case "client":
		username := r.FormValue("username")
		fullname := r.FormValue("fullname")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		password := r.FormValue("password")

		var role models.Role
		if err := database.DB.Model(&models.Role{}).Where("name = ?", "client").Find(&role).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}

		user := models.User{
			Fullname:    fullname,
			PhoneNumber: phone,
			Email:       email,
			Password:    password,
			Username:    username,
			RoleID:      role.ID,
		}
		database.DB.Create(&user)

		helpers.SetSession(w, user.ID)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /log_out
// Method: POST
func LogOut(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	if err := database.DB.Where("token_hash = ?", tokenHash).Delete(models.Session{}).Error; err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
