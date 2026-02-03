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
	var role string
	if err == nil {
		token := sessionToken.Value
		user = helpers.GetUserBySessionToken(token)
		role = user.Role.Name
	}

	query := r.URL.Query()
	gateType := query.Get("gateType")

	cfg, err := getGateCfg(gateType)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	isDealer := (role == "dealer" || role == "admin" || role == "manager")

	data := map[string]any{
		"css":      "calculator.css",
		"cfg":      cfg,
		"user":     user,
		"isDealer": isDealer,
	}

	if err := templates.ExecuteTemplate(w, "calc.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
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
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"price": size.WholesalePrice,
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
			helpers.WriteError(w, err, http.StatusInternalServerError)
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
		WidthParams:  types.SizeParams{},
		HeightParams: types.SizeParams{},
		Options:      []models.Option{},
		Products:     []models.Product{},
	}

	var gateTypes []models.GateType
	if err := database.DB.Find(&gateTypes).Error; err != nil {
		return types.Config{}, err
	}

	gateTypesMap := make(map[int32]string)
	for _, gateTypeItem := range gateTypes {
		gateTypesMap[gateTypeItem.ID] = gateTypeItem.Name
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	errChan := make(chan error, 8)
	gateTypeInt, err := strconv.ParseInt(gateType, 10, 32)
	if err != nil {
		return cfg, err
	}

	switch gateTypesMap[int32(gateTypeInt)] {
	case "Промышленные ворота":
		wg.Go(func() {
			var drives []models.IndustrialGateDrive
			if err := database.DB.Find(&drives).Error; err != nil {
				errChan <- err
				return
			}
			mutex.Lock()
			cfg.IndustrialDrives = drives
			mutex.Unlock()
		})
	case "Бытовые ворота":
		wg.Go(func() {
			var drives []models.ResidentialGateDrive
			if err := database.DB.Find(&drives).Error; err != nil {
				errChan <- err
				return
			}
			mutex.Lock()
			cfg.ResidentialDrives = drives
			mutex.Unlock()
		})
	}

	wg.Go(func() {
		var liftTypes []models.LiftType
		if err := database.DB.Find(&liftTypes).Error; err != nil {
			errChan <- err
			return
		}
		mutex.Lock()
		cfg.LiftTypes = liftTypes
		mutex.Unlock()
	})

	wg.Go(func() {
		var colors []models.Color
		if err := database.DB.Find(&colors).Error; err != nil {
			errChan <- err
			return
		}
		mutex.Lock()
		cfg.Colors = colors
		mutex.Unlock()
	})

	wg.Go(func() {
		var amounts []models.CycleAmount
		if err := database.DB.Find(&amounts).Error; err != nil {
			errChan <- err
			return
		}
		mutex.Lock()
		cfg.CycleAmounts = amounts
		mutex.Unlock()
	})

	wg.Go(func() {
		var params types.SizeParams
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(width) as max_value, MIN(width) as min_value").
			Where("gate_type_id = ?", gateType).
			Scan(&params).Error; err != nil {
			errChan <- err
			return
		}

		mutex.Lock()
		cfg.WidthParams = params
		mutex.Unlock()
	})

	wg.Go(func() {
		var params types.SizeParams
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(height) as max_value, MIN(height) as min_value").
			Where("gate_type_id = ?", gateType).
			Scan(&params).Error; err != nil {
			errChan <- err
			return
		}

		mutex.Lock()
		cfg.HeightParams = params
		mutex.Unlock()
	})

	wg.Go(func() {
		var products []models.Product
		if err := database.DB.Find(&products).Error; err != nil {
			errChan <- err
			return
		}
		mutex.Lock()
		cfg.Products = products
		mutex.Unlock()
	})

	wg.Go(func() {
		var options []models.Option
		if err := database.DB.Find(&options).Error; err != nil {
			errChan <- err
			return
		}
		mutex.Lock()
		cfg.Options = options
		mutex.Unlock()
	})

	wg.Wait()
	close(errChan)

	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}

	return cfg, firstErr
}
