package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/rails"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/samborkent/uuidv7"
)

func UploadProductImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	dir := os.Getenv("IMAGES_DIRECTORY")
	if dir == "" {
		return "", fmt.Errorf("переменная IMAGES_DIRECTORY не задана")
	}

	ext := filepath.Ext(handler.Filename)
	filename := uuidv7.New().String() + ext

	path := filepath.Join(dir, filename)
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}

func AddNewDataBaseTableRow(tableName string, tableData any) error {
	switch tableName {
	case "colors":
		tableData, ok := tableData.(colors.Color)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewColor(tableData); err != nil {
			return err
		}
	case "cycle_amounts":
		tableData, ok := tableData.(cycle_amounts.CycleAmount)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewCycleAmount(tableData); err != nil {
			return err
		}
	case "industrial_drives":
		tableData, ok := tableData.(industrial_gate_drives.IndustrialGateDrive)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewIndustrialDrive(tableData); err != nil {
			return err
		}
	case "lift_types":
		tableData, ok := tableData.(lift_types.LiftType)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewLiftType(tableData); err != nil {
			return err
		}
	case "options":
		tableData, ok := tableData.(options.Option)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewOption(tableData); err != nil {
			return err
		}
	case "products":
		tableData, ok := tableData.(products.Product)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewProduct(tableData); err != nil {
			return err
		}
	case "rails":
		tableData, ok := tableData.(rails.Rail)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewRail(tableData); err != nil {
			return err
		}
	case "residential_drives":
		tableData, ok := tableData.(residential_gate_drives.ResidentialGateDrive)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.CreateNewResidentialDrive(tableData); err != nil {
			return err
		}
	}

	return nil
}

func GetTablePageData(tableName string) (any, error) {
	switch tableName {
	case "colors":
		return repository.GetColors()
	case "cycle_amounts":
		return repository.GetCycleAmounts()
	case "industrial_drives":
		return repository.GetIndustrialDrives()
	case "lift_types":
		return repository.GetLiftTypes()
	case "options":
		return repository.GetOptions()
	case "products":
		return repository.GetProducts()
	case "rails":
		return repository.GetRails()
	case "residential_drives":
		return repository.GetResidentialDrives()
	case "users":
		return repository.GetAllUsers()
	default:
		return nil, errs.ErrInvalidTableType
	}
}

func UpdateRow(tableName string, tableData any) error {
	switch tableName {
	case "colors":
		tableData, ok := tableData.(colors.Color)
		if !ok {
			return errs.ErrInvalidTableType
		}

		return repository.UpdateColor(tableData)
	case "cycle_amounts":
		tableData, ok := tableData.(cycle_amounts.CycleAmount)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateCycleAmount(tableData); err != nil {
			return err
		}
	case "industrial_drives":
		tableData, ok := tableData.(industrial_gate_drives.IndustrialGateDrive)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateIndustrialDrive(tableData); err != nil {
			return err
		}
	case "lift_types":
		tableData, ok := tableData.(lift_types.LiftType)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateLiftType(tableData); err != nil {
			return err
		}
	case "options":
		tableData, ok := tableData.(options.Option)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateOption(tableData); err != nil {
			return err
		}
	case "products":
		tableData, ok := tableData.(products.Product)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateProduct(tableData); err != nil {
			return err
		}
	case "rails":
		tableData, ok := tableData.(rails.Rail)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateRail(tableData); err != nil {
			return err
		}
	case "residential_drives":
		tableData, ok := tableData.(residential_gate_drives.ResidentialGateDrive)
		if !ok {
			return errs.ErrInvalidTableType
		}

		if err := repository.UpdateResidentialDrive(tableData); err != nil {
			return err
		}
	}

	return nil
}

func DeleteRow(tableName string, rowId int64) error {
	switch tableName {
	case "colors":
		return repository.DeleteColor(rowId)
	case "cycle_amounts":
		return repository.DeleteCycleAmount(rowId)
	case "industrial_drives":
		return repository.DeleteIndustrialDrive(rowId)
	case "lift_types":
		return repository.DeleteLiftType(rowId)
	case "options":
		return repository.DeleteOption(rowId)
	case "products":
		return repository.DeleteProduct(rowId)
	case "rails":
		return repository.DeleteRail(rowId)
	case "residential_drives":
		return repository.DeleteResidentialDrive(rowId)
	case "users":
		return repository.AdminDeleteUser(rowId)
	}

	return nil
}
