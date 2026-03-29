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

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

func hasDealerAccess(role string) bool {
	return role == "dealer" || role == "admin" || role == "manager"
}

// Route: /calculator
// Method: GET
func GetCalculatorForUser(w http.ResponseWriter, r *http.Request) {
	user, _ := helpers.GetUserBySessionToken(r)

	role := "client"
	if user != nil {
		role = user.Role.Name
	}

	isDealer := hasDealerAccess(role)

	indCfg, err := getGateCfg(models.GateTypeInd, isDealer)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	resCfg, err := getGateCfg(models.GateTypeRes, isDealer)
	if err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":                      "calc.css",
		"IndustrialConfiguration":  indCfg,
		"ResidentialConfiguration": resCfg,
		"user":                     user,
		"isDealer":                 isDealer,
	}

	if err := templates.ExecuteTemplate(w, "calc.html", data); err != nil {
		helpers.WriteError(w, err, http.StatusInternalServerError)
	}
}

// Route: /sizes
// Method: GET
func GetPriceBasedOnSize(w http.ResponseWriter, r *http.Request) {
	user, _ := helpers.GetUserBySessionToken(r)

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

	gateType, err := models.DetermineGateType(query.Get("gateType"))
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
			Where("gate_type = ?", gateType).
			Limit(1).
			Order("width asc, height asc").
			Find(&size).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		resp = &proto_files.SizePrice{
			Price: &proto_files.SizePrice_Dealer{
				Dealer: &proto_files.DealerSizePrices{
					DealerPrice: size.WholesalePrice.Mul(decimal.NewFromInt(100)).IntPart(),
					ClientPrice: size.RetailPrice.Mul(decimal.NewFromInt(100)).IntPart(),
				},
			},
		}
	} else {
		var price decimal.Decimal
		if err := database.DB.Model(&models.Size{}).
			Select("RetailPrice").
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type = ?", gateType).
			Limit(1).
			Order("width asc, height asc").
			Pluck("RetailPrice", &price).Error; err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		resp = &proto_files.SizePrice{
			Price: &proto_files.SizePrice_Client{
				Client: &proto_files.ClientSizePrice{
					ClientPrice: price.Mul(decimal.NewFromInt(100)).IntPart(),
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

func getGateCfg(gateType models.GateType, isDealer bool) (types.Config, error) {
	cfg := types.Config{
		LiftTypes:    []models.LiftType{},
		Colors:       []models.Color{},
		WidthParams:  types.SizeParams{},
		HeightParams: types.SizeParams{},
		Options:      []models.Option{},
		Products:     []models.Product{},
		DriveTypes:   []models.DriveType{},
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 8)

	switch gateType {
	case models.GateTypeInd:
		cfg.DriveTypes = []models.DriveType{models.IndDriveType, models.ManualDriveType}
		wg.Go(func() {
			if err := database.DB.Find(&cfg.IndustrialDrives).Error; err != nil {
				errChan <- err
				return
			}
		})

		if !isDealer {
			for i := range cfg.IndustrialDrives {
				cfg.IndustrialDrives[i].WholesalePrice = decimal.NewFromInt(0)
			}
		}
	case models.GateTypeRes:
		cfg.DriveTypes = []models.DriveType{models.ResDriveType, models.ManualDriveType}
		wg.Go(func() {
			if err := database.DB.Find(&cfg.Rails).Error; err != nil {
				errChan <- err
				return
			}
		})

		if !isDealer {
			for i := range cfg.Rails {
				cfg.Rails[i].WholesalePrice = decimal.NewFromInt(0)
			}
		}

		wg.Go(func() {
			if err := database.DB.Find(&cfg.ResidentialDrives).Error; err != nil {
				errChan <- err
				return
			}
		})

		if !isDealer {
			for i := range cfg.ResidentialDrives {
				cfg.ResidentialDrives[i].WholesalePrice = decimal.NewFromInt(0)
			}
		}
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
			Where("gate_type = ?", gateType).
			Scan(&cfg.WidthParams).Error; err != nil {
			errChan <- err
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(height) as max_value, MIN(height) as min_value").
			Where("gate_type = ?", gateType).
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

	if !isDealer {
		for i := range cfg.LiftTypes {
			cfg.LiftTypes[i].WholesaleMarkup = decimal.NewFromInt(0)
		}
		for i := range cfg.CycleAmounts {
			cfg.CycleAmounts[i].WholesaleMarkup = decimal.NewFromInt(0)
		}
		for i := range cfg.Products {
			cfg.Products[i].WholesalePrice = decimal.NewFromInt(0)
		}
		for i := range cfg.Options {
			cfg.Options[i].WholesalePrice = decimal.NewFromInt(0)
		}
	}

	return cfg, <-errChan
}
