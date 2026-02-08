package handlers

import (
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	//"github.com/gorilla/mux"
)

// Route: /tables/{table_name}
// Method: GET
func GetDataBaseRedactor(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)
	if user == nil {
		return
	}

	role := user.Role.Name
	if role != "admin" {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}
	//vars := mux.Vars(r)

	//model := database.NamedModels[vars["table_name"]]

}

// Route: /tables
// Method: GET
func GetDataBaseTableList(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)
	if user == nil {
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
