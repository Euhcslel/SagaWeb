package http

import (
	"net/http"
	"project/internal/database"
	"project/internal/helpers"
)

// Route: /tables/{table_name}
// Method: GET
func GetDataBaseRedactor(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	role := user.Role.Name
	if role != "admin" {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}
}

// Route: /tables
// Method: GET
func GetDataBaseTableList(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetUserBySessionToken(r)
	if err != nil {
		helpers.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	role := user.Role.Name
	if role != "admin" {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	migrator := database.DB.Migrator()
	tables, err := migrator.GetTables()
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}

	data := map[string]any{
		"css":    "",
		"tables": tables,
		"user":   user,
	}

	if err := templates.ExecuteTemplate(w, "tables.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}
