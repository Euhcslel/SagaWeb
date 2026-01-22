package handlers

import (
	"html/template"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
)

var templates = template.Must(template.ParseGlob("templates/**/*.html"))
var cssPath = "/assets/styles/"

// Route: /
// Method: GET
func MainHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		user = models.User{}
	} else {
		token := sessionToken.Value
		user = helpers.GetUserBySessionToken(token)
	}

	data := map[string]any{
		"css":  cssPath+"main.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "main.html", data); err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
	}
}

// Route: /contacts
// Method: GET
func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "contacts.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /gate_types
// Method: GET
func GetGateTypesList(w http.ResponseWriter, r *http.Request) {
	var gateTypes []models.GateType
	err := database.DB.Find(&gateTypes).Error
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":       "",
		"gateTypes": gateTypes,
	}

	if err := templates.ExecuteTemplate(w, "gate_type.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
	}
}
