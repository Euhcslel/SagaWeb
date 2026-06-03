package http

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
	"github.com/shopspring/decimal"
)

// Функция для форматирования даты и времени в шаблонах
func formatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

// Функция для вывода цены в корректном формате
func formatPrice(d decimal.Decimal) string {
	d = d.Round(2)
	intPart := d.Truncate(0).String()
	fracPart := d.Sub(d.Truncate(0)).Mul(decimal.NewFromInt(100)).IntPart()

	// Форматируем целую часть
	var result strings.Builder
	n := len(intPart)
	for i, ch := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune(' ')
		}
		result.WriteRune(ch)
	}

	return fmt.Sprintf("%s,%02d ₽", result.String(), fracPart)
}

var templates = template.Must(
	template.New("").
		Funcs(template.FuncMap{
			"fmtTime":     formatDateTime,
			"fmtPrice": formatPrice,
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
		"css":     "home.css",
		"user":    user,
	}

	if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}
