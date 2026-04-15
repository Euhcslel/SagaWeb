package http

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/service"
	"project/internal/utils"
	"time"
)

// Route: /sign_in
// Method: GET
func SignInForm(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())
	if user != nil {
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
	currentUser := utils.UserFromContext(r.Context())
	if currentUser != nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userId, err := service.SignIn(username, password)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	if err = utils.SetSession(w, userId); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /sign_up
// Method: GET
func SignUpForm(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())
	if user != nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"css": "auth.css",
	}
	if err := templates.ExecuteTemplate(w, "sign_up.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /sign_up
// Method: POST
func SignUp(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())
	if user != nil {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	fullname := r.FormValue("fullname")
	company := r.FormValue("company")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err = service.SignUp(fullname, company, phone, email, password)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /sign_out
// Method: POST
func SignOut(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session_token"); err == nil {
		if err := service.SignOut(cookie.Value); err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
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
