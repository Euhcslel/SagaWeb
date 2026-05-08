package http

import (
	"html/template"
	"net/http"
	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"time"
)

// Функция для форматирования даты и времени в шаблонах
func FormatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

var templates = template.Must(
	template.New("").
		Funcs(template.FuncMap{
			"fmtTime": FormatDateTime,
		}).
		ParseGlob("web/templates/**/*.html"),
)

// Route: /
// Method: GET
func MainHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	dealers, err := service.GetDealersList()
	if err != nil {
		dealers = nil
	}

	data := map[string]any{
		"dealers": dealers,
		"css":  "home.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}
