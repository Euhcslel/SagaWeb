package service

import (
	"errors"

	"github.com/Euhcslel/SagaWeb/internal/domain/colors"
	"github.com/Euhcslel/SagaWeb/internal/domain/cycle_amount"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/lift_types"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/products"
	"github.com/Euhcslel/SagaWeb/internal/domain/residential_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/repository"
)

func AddNewDataBaseTableRow(tableName string, tableData any) error {
	switch tableName {
	case "colors":
		tableData, ok := tableData.(colors.Color)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewColor(tableData); err != nil {
			return err
		}
	case "cycle_amounts":
		tableData, ok := tableData.(cycle_amount.CycleAmount)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewCycleAmount(tableData); err != nil {
			return err
		}
	case "industrial_drives":
		tableData, ok := tableData.(industrial_gate_drives.IndustrialGateDrive)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewIndustrialDrive(tableData); err != nil {
			return err
		}
	case "lift_types":
		tableData, ok := tableData.(lift_types.LiftType)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewLiftType(tableData); err != nil {
			return err
		}
	case "options":
		tableData, ok := tableData.(options.Option)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewOption(tableData); err != nil {
			return err
		}
	case "products":
		tableData, ok := tableData.(products.Product)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewProduct(tableData); err != nil {
			return err
		}
	case "residential_drives":
		tableData, ok := tableData.(residential_gate_drives.ResidentialGateDrive)
		if !ok {
			return errors.New("invalid table type")
		}

		if err := repository.CreateNewResidentialDrive(tableData); err != nil {
			return err
		}
	}

	return nil
}
