package service

import (
	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amounts"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"
)

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
	case "residential_drives":
		return repository.GetResidentialDrives()
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
	case "residential_drives":
		return repository.DeleteResidentialDrive(rowId)
	}

	return nil
}
