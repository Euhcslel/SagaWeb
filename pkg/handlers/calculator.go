package handlers

import (
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"project/pkg/proto_files"
	"project/pkg/types"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"
)

func hasDealerAccess(role string) bool {
	return role == "dealer" || role == "admin" || role == "manager"
}

// Route: /calculator
// Method: GET
func GetCalculatorForUser(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserBySessionToken(w, r)
	role := "client"
	if user != nil {
		role = user.Role.Name
	}

	query := r.URL.Query()
	gateType := query.Get("gateType")

	cfg, err := getGateCfg(gateType)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	isDealer := hasDealerAccess(role)

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
	user := helpers.GetUserBySessionToken(w, r)
	role := "client"
	if user != nil {
		role = user.Role.Name
	}

	query := r.URL.Query()
	width, err := strconv.Atoi(query.Get("width"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(query.Get("height"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	gateTypeID, err := strconv.Atoi(query.Get("gateType"))
	if err != nil {
		helpers.WriteError(w, err, http.StatusBadRequest)
		return
	}

	var data []byte
	var resp *proto_files.SizePrice

	if hasDealerAccess(role) {
		var size models.Size
		if err := database.DB.Model(&models.Size{}).
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type_id = ?", gateTypeID).
			Limit(1).
			Order("width asc, height asc").
			Find(&size).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		resp = &proto_files.SizePrice{
			Price: &proto_files.SizePrice_Dealer{
				Dealer: &proto_files.DealerSizePrices{
					DealerPrice: size.WholesalePrice,
					ClientPrice: size.RetailPrice,
				},
			},
		}
	} else {
		var price int64
		if err := database.DB.Model(&models.Size{}).
			Select("RetailPrice").
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type_id = ?", gateTypeID).
			Limit(1).
			Order("width asc, height asc").
			Pluck("RetailPrice", &price).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		resp = &proto_files.SizePrice{
			Price: &proto_files.SizePrice_Client{
				Client: &proto_files.ClientSizePrice{
					ClientPrice: price,
				},
			},
		}

	}

	data, err = proto.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
	errChan := make(chan error, 8)
	gateTypeInt, err := strconv.ParseInt(gateType, 10, 32)
	if err != nil {
		return cfg, err
	}

	switch gateTypesMap[int32(gateTypeInt)] {
	case "Промышленные ворота":
		wg.Go(func() {
			if err := database.DB.Find(&cfg.IndustrialDrives).Error; err != nil {
				errChan <- err
				return
			}
		})
	case "Бытовые ворота":
		wg.Go(func() {
			if err := database.DB.Find(&cfg.ResidentialDrives).Error; err != nil {
				errChan <- err
				return
			}
		})
	}

	wg.Go(func() {
		if err := database.DB.Find(&cfg.LiftTypes).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Colors).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.CycleAmounts).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(width) as max_value, MIN(width) as min_value").
			Where("gate_type_id = ?", gateType).
			Scan(&cfg.WidthParams).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(height) as max_value, MIN(height) as min_value").
			Where("gate_type_id = ?", gateType).
			Scan(&cfg.HeightParams).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Products).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Options).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Wait()
	close(errChan)

	return cfg, <-errChan
}
