package service

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	"github.com/Euhcslel/SagaWeb/internal/generated"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/Euhcslel/SagaWeb/internal/utils"

	"github.com/shopspring/decimal"
)

// calculatorPageData описывает данные необходимые для страницы калькулятора.
type calculatorPageData struct {
	IndustrialConfiguration  any
	ResidentialConfiguration any
	IsDealer                 bool
}

// GetCalculatorPageData возвращает данные для страницы калькулятора.
func GetCalculatorPageData(user *users.User) (*calculatorPageData, error) {
	role := enums.ClientRole
	if user != nil {
		role = user.Role
	}

	isDealer := utils.HasDealerAccess(role)

	indCfg, err := GetGateCfg(enums.GateTypeInd, isDealer)
	if err != nil {
		return nil, err
	}

	resCfg, err := GetGateCfg(enums.GateTypeRes, isDealer)
	if err != nil {
		return nil, err
	}

	return &calculatorPageData{
		IndustrialConfiguration:  indCfg,
		ResidentialConfiguration: resCfg,
		IsDealer:                 isDealer,
	}, nil
}

// GetPriceBasedOnSize возвращает цену на ворота в зависимости от размеров и типа ворот.
func GetPriceBasedOnSize(width, height int64, gateType enums.GateType, user *users.User) (*generated.SizePrice, error) {
	role := enums.ClientRole
	if user != nil {
		role = user.Role
	}

	var resp *generated.SizePrice

	size, err := repository.GetSizeForDimensions(width, height, gateType)
	if err != nil {
		return nil, err
	}

	if utils.HasDealerAccess(role) {
		resp = &generated.SizePrice{
			Price: &generated.SizePrice_Dealer{
				Dealer: &generated.DealerSizePrices{
					DealerPrice: size.WholesalePrice.Mul(decimal.NewFromInt(100)).IntPart(),
					ClientPrice: size.RetailPrice.Mul(decimal.NewFromInt(100)).IntPart(),
				},
			},
		}
	} else {
		resp = &generated.SizePrice{
			Price: &generated.SizePrice_Client{
				Client: &generated.ClientSizePrice{
					ClientPrice: size.RetailPrice.Mul(decimal.NewFromInt(100)).IntPart(),
				},
			},
		}
	}

	return resp, nil
}

// GetGateCfg возвращает конфигурацию ворот в зависимости от типа ворот и роли пользователя.
func GetGateCfg(gateType enums.GateType, isDealer bool) (*types.Config, error) {
	cfg := &types.Config{
		LiftTypes:    []lift_types.LiftType{},
		Colors:       []colors.Color{},
		WidthParams:  &types.SizeParams{},
		HeightParams: &types.SizeParams{},
		Options:      []options.Option{},
		Products:     []products.Product{},
		DriveTypes:   []enums.DriveType{},
		CycleAmounts: []cycle_amounts.CycleAmount{},
	}

	var err error

	switch gateType {
	case enums.GateTypeInd:
		cfg.DriveTypes = []enums.DriveType{enums.IndDriveType, enums.ManualDriveType}

		cfg.IndustrialDrives, err = repository.GetIndustrialDrives()
		if err != nil {
			return nil, err
		}

		if !isDealer {
			for i := range cfg.IndustrialDrives {
				cfg.IndustrialDrives[i].WholesalePrice = decimal.NewFromInt(0)
			}
		}
	case enums.GateTypeRes:
		cfg.DriveTypes = []enums.DriveType{enums.ResDriveType, enums.ManualDriveType}

		cfg.Rails, err = repository.GetRails()
		if err != nil {
			return nil, err
		}

		if !isDealer {
			for i := range cfg.Rails {
				cfg.Rails[i].WholesalePrice = decimal.NewFromInt(0)
			}
		}

		cfg.ResidentialDrives, err = repository.GetResidentialDrives()
		if err != nil {
			return nil, err
		}

		if !isDealer {
			for i := range cfg.ResidentialDrives {
				cfg.ResidentialDrives[i].WholesalePrice = decimal.NewFromInt(0)
			}
		}
	}

	cfg.ManualDrive, err = repository.GetManualDrivePrices()
	if err != nil {
		return nil, err
	}

	cfg.LiftTypes, err = repository.GetLiftTypes()
	if err != nil {
		return nil, err
	}

	cfg.Colors, err = repository.GetColors()
	if err != nil {
		return nil, err
	}

	cfg.CycleAmounts, err = repository.GetCycleAmounts()
	if err != nil {
		return nil, err
	}

	cfg.WidthParams, err = repository.GetMaxAndMinWidth(gateType)
	if err != nil {
		return nil, err
	}

	cfg.HeightParams, err = repository.GetMaxAndMinHeight(gateType)
	if err != nil {
		return nil, err
	}

	cfg.Products, err = repository.GetProducts()
	if err != nil {
		return nil, err
	}

	cfg.Options, err = repository.GetOptions()
	if err != nil {
		return nil, err
	}

	if !isDealer {
		cfg.ManualDrive.ChainMeterWholesalePrice = decimal.NewFromInt(0)
		cfg.ManualDrive.RcpWholesalePrice = decimal.NewFromInt(0)

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

	return cfg, nil
}
