package service

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/utils"
)

// CatalogPageData описывает информацию для страницы каталога
type CatalogPageData struct {
	Products          []products.Product
	Options           []options.Option
	IndustrialDrives  []industrial_gate_drives.IndustrialGateDrive
	ResidentialDrives []residential_gate_drives.ResidentialGateDrive
	Rails             []rails.Rail
	IsDealer          bool
}

// GetCatalogPageData возвращает данные для страницы каталога.
func GetCatalogPageData(user *users.User) (*CatalogPageData, error) {
	role := enums.ClientRole
	if user != nil {
		role = user.Role
	}
	isDealer := utils.HasDealerAccess(role)

	prods, err := repository.GetProducts()
	if err != nil {
		return nil, err
	}

	opts, err := repository.GetOptions()
	if err != nil {
		return nil, err
	}

	indDrives, err := repository.GetIndustrialDrives()
	if err != nil {
		return nil, err
	}

	resDrives, err := repository.GetResidentialDrives()
	if err != nil {
		return nil, err
	}

	railsList, err := repository.GetRails()
	if err != nil {
		return nil, err
	}

	return &CatalogPageData{
		Products:          prods,
		Options:           opts,
		IndustrialDrives:  indDrives,
		ResidentialDrives: resDrives,
		Rails:             railsList,
		IsDealer:          isDealer,
	}, nil
}
