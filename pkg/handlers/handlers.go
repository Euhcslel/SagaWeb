package handlers

import (
	"html/template"
	"net/http"
	"project/pkg/helpers"
	"project/pkg/models"
	"time"
)

// Функция, для создания словаря из пар ключ-значение (для html-шаблонов)
func dict(values ...any) map[string]any {
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}

// Функция для форматирования даты и времени в шаблонах
func FormatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

var templates = template.Must(
	template.New("").
		Funcs(template.FuncMap{
			"dict": dict,
			"fmtTime": FormatDateTime,
		}).
		ParseGlob("templates/**/*.html"),
)

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

	data := map[string]any{
		"css":       "",
		"gateTypes": models.GetAllGateTypes(),
		"user":      user,
	}

	if err := templates.ExecuteTemplate(w, "gate_type.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}
