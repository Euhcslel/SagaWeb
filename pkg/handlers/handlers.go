package handlers

import (
	"html/template"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
)

var templates = template.Must(template.ParseGlob("templates/**/*.html"))

// Route: /
// Method: GET
func MainHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)

	data := map[string]any{
		"css":  "main.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "main.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /gate_types
// Method: GET
func GetGateTypesList(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)

	var gateTypes []models.GateType
	err := database.DB.Find(&gateTypes).Error
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":       "",
		"gateTypes": gateTypes,
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "gate_type.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}
