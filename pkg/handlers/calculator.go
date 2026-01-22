package handlers

import (
	"encoding/json"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"project/pkg/types"
	"strconv"
	"sync"
)

// Route: /calculator
// Method: GET
func GetCalculatorForUser(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	var user models.User
	if err == nil {
		token := sessionToken.Value
		user = helpers.GetUserBySessionToken(token)
	}

	query := r.URL.Query()
	gateType := query.Get("gateType")

	cfg, err := getGateCfg(gateType)
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":  "/assets/styles/calculator.css",
		"cfg":  cfg,
		"user": user,
	}

	role := user.Role.Name
	if role == "dealer" || role == "admin" || role == "manager" {
		if err := templates.ExecuteTemplate(w, "dealer_calc.html", data); err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		}
	} else {
		if err := templates.ExecuteTemplate(w, "client_calc.html", data); err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		}
	}
}

// Route: /sizes
// Method: GET
func GetPriceBasedOnSize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width := query.Get("width")
	height := query.Get("height")
	gateType := query.Get("gateType")

	sessionToken, err := r.Cookie("session_token")
	var role string
	if err == nil {
		token := sessionToken.Value
		user := helpers.GetUserBySessionToken(token)
		role = user.Role.Name
	}

	if role == "dealer" || role == "manager" || role == "admin" {
		var size models.Size
		if err := database.DB.Model(&models.Size{}).
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type_id = ?", gateType).
			Limit(1).
			Order("width asc, height asc").
			Find(&size).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"dealer_price": size.WholesalePrice,
			"client_price": size.RetailPrice,
		})
	} else {
		var price int64
		if err := database.DB.Model(&models.Size{}).
			Select("RetailPrice").
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type_id = ?", gateType).
			Limit(1).
			Order("width asc, height asc").
			Pluck("RetailPrice", &price).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"price": price,
		})
	}

}

func getGateCfg(gateType string) (types.Config, error) {
	cfg := types.Config{
		LiftTypes:    []models.LiftType{},
		Colors:       []models.Color{},
		WidthParams:  types.WidthParams{},
		HeightParams: types.HeightParams{},
		Options:      []models.Option{},
		Products:     []models.Product{},
	}

	var returnErr error
	var gateTypes []models.GateType
	if err := database.DB.Find(&gateTypes).Error; err != nil {
		return types.Config{}, err
	}

	gateTypesMap := make(map[int32]string)
	for _, gateTypeItem := range gateTypes {
		gateTypesMap[gateTypeItem.ID] = gateTypeItem.Name
	}

	var wg sync.WaitGroup
	gateTypeInt, _ := strconv.ParseInt(gateType, 10, 32)
	switch gateTypesMap[int32(gateTypeInt)] {
	case "Промышленные ворота":
		wg.Go(func() {
			cfg.IndustrialDrives = []models.IndustrialGateDrive{}
			if err := database.DB.Find(&cfg.IndustrialDrives).Error; err != nil {
				returnErr = err
			}
		})
	case "Бытовые ворота":
		wg.Go(func() {
			cfg.ResidentialDrives = []models.ResidentialGateDrive{}
			if err := database.DB.Find(&cfg.ResidentialDrives).Error; err != nil {
				returnErr = err
			}
		})
	}

	wg.Go(func() {
		if err := database.DB.Find(&cfg.LiftTypes).Error; err != nil {
			returnErr = err
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Colors).Error; err != nil {
			returnErr = err
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.CycleAmounts).Error; err != nil {
			returnErr = err
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(width) as max_width, MIN(width) as min_width").
			Where("gate_type_id = ?", gateType).
			Scan(&cfg.WidthParams).Error; err != nil {
			returnErr = err
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(height) as max_height, MIN(height) as min_height").
			Where("gate_type_id = ?", gateType).
			Scan(&cfg.HeightParams).Error; err != nil {
			returnErr = err
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Products).Error; err != nil {
			returnErr = err
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Options).Error; err != nil {
			returnErr = err
		}
	})

	wg.Wait()

	return cfg, returnErr
}
