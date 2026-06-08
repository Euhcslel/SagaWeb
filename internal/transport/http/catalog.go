package http

import (
	"net/http"

	"github.com/Euhcslel/SagaWeb/internal/helpers"
	"github.com/Euhcslel/SagaWeb/internal/service"
	"github.com/Euhcslel/SagaWeb/internal/utils"
)

// Route: /catalog
// Method: GET
func GetCatalog(w http.ResponseWriter, r *http.Request) {
	user := utils.UserFromContext(r.Context())

	pageData, err := service.GetCatalogPageData(user)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":               "catalog.css",
		"user":              user,
		"isDealer":          pageData.IsDealer,
		"Products":          pageData.Products,
		"Options":           pageData.Options,
		"IndustrialDrives":  pageData.IndustrialDrives,
		"ResidentialDrives": pageData.ResidentialDrives,
		"Rails":             pageData.Rails,
	}

	if err := templates.ExecuteTemplate(w, "catalog.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}
