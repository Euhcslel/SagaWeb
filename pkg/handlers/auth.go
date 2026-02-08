package handlers

import (
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Route: /sign_in
// Method: GET
func SignInForm(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"css": "auth.css",
	}

	if err := templates.ExecuteTemplate(w, "sign_in.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /sign_in
// Method: POST
func SignIn(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
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
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	userId := user.ID
	dbPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	helpers.SetSession(w, userId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /sign_up
// Method: GET
func SignUpFrom(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"css": "auth.css",
	}
	if err := templates.ExecuteTemplate(w, "sign_up.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /sign_up
// Method: POST
func SignUp(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	fullname := r.FormValue("fullname")
	company := r.FormValue("company")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	password := r.FormValue("password")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var status models.Status
	if err := database.DB.Model(&models.Status{}).Where("name = ?", "Ожидает подтверждения").Find(&status).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	request := models.DealerRegRequest{
		Company:     company,
		Fullname:    fullname,
		PhoneNumber: phone,
		Email:       email,
		Password:    passwordHash,
		StatusID:    status.ID,
	}
	database.DB.Create(&request)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /sign_out
// Method: POST
func SignOut(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value

	if err := database.DB.Where("token = ?", token).Delete(models.Session{}).Error; err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
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
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
