package service

import (
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)


func UpdateOrderByLogistician(user *users.User, orderID int64, manufactureDate *time.Time) error {
	if user.Role != enums.LogisticianRole {
		return errs.ErrForbidden
	}
	return repository.UpdateOrderManufactureDate(orderID, manufactureDate)
}

func ImportSizePrices(
	indDealerFile multipart.File,
	indClientFile multipart.File,
	resDealerFile multipart.File,
	resClientFile multipart.File,
) error {
	if err := updateSizes(indDealerFile, indClientFile, enums.GateTypeInd); err != nil {
		return fmt.Errorf("промышленные ворота: %w", err)
	}
	if err := updateSizes(resDealerFile, resClientFile, enums.GateTypeRes); err != nil {
		return fmt.Errorf("бытовые ворота: %w", err)
	}
	return nil
}

func updateSizes(dealerFile multipart.File, clientFile multipart.File, gateType enums.GateType) error {
	clientXlsx, err := excelize.OpenReader(clientFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла клиентских цен: %w", err)
	}
	defer clientXlsx.Close()

	dealerXlsx, err := excelize.OpenReader(dealerFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла дилерских цен: %w", err)
	}
	defer dealerXlsx.Close()

	clientRows, err := clientXlsx.GetRows("Лист1")
	if err != nil {
		return fmt.Errorf("ошибка чтения строк клиентского файла: %w", err)
	}
	clientCols, err := clientXlsx.GetCols("Лист1")
	if err != nil {
		return fmt.Errorf("ошибка чтения столбцов клиентского файла: %w", err)
	}
	dealerRows, err := dealerXlsx.GetRows("Лист1")
	if err != nil {
		return fmt.Errorf("ошибка чтения строк дилерского файла: %w", err)
	}
	dealerCols, err := dealerXlsx.GetCols("Лист1")
	if err != nil {
		return fmt.Errorf("ошибка чтения столбцов дилерского файла: %w", err)
	}

	if len(dealerCols) != len(clientCols) || len(dealerRows) != len(clientRows) {
		return fmt.Errorf("число столбцов или строк в файлах различается")
	}

	// Логика идентична cmd/size_import: i и j — 1-индексированные адреса Excel
	for i := 1; i < len(clientRows); i++ {
		heightStr, err := clientXlsx.GetCellValue("Лист1", "A"+strconv.Itoa(i))
		if err != nil {
			return fmt.Errorf("ошибка чтения высоты (A%d): %w", i, err)
		}
		if heightStr == "" {
			continue
		}
		height, err := strconv.ParseInt(heightStr, 10, 64)
		if err != nil {
			return fmt.Errorf("невалидная высота %q: %w", heightStr, err)
		}

		for j := 1; j < len(clientCols); j++ {
			if j == 1 || i == 1 {
				continue
			}
			colLetter, err := excelize.ColumnNumberToName(j)
			if err != nil {
				return fmt.Errorf("ошибка преобразования номера столбца: %w", err)
			}

			widthStr, err := clientXlsx.GetCellValue("Лист1", colLetter+"1")
			if err != nil {
				return fmt.Errorf("ошибка чтения ширины (%s1): %w", colLetter, err)
			}
			if widthStr == "" {
				continue
			}
			width, err := strconv.ParseInt(widthStr, 10, 64)
			if err != nil {
				return fmt.Errorf("невалидная ширина %q: %w", widthStr, err)
			}

			clientCell, err := clientXlsx.GetCellValue("Лист1", colLetter+strconv.Itoa(i))
			if err != nil {
				return fmt.Errorf("ошибка чтения клиентской цены (%s%d): %w", colLetter, i, err)
			}
			dealerCell, err := dealerXlsx.GetCellValue("Лист1", colLetter+strconv.Itoa(i))
			if err != nil {
				return fmt.Errorf("ошибка чтения дилерской цены (%s%d): %w", colLetter, i, err)
			}

			if clientCell == "" || dealerCell == "" {
				log.Printf("пропуск: ширина %s, высота %s — нет значения\n", widthStr, heightStr)
				continue
			}

			retailRaw, err := strconv.Atoi(clientCell)
			if err != nil {
				return fmt.Errorf("невалидная клиентская цена %q: %w", clientCell, err)
			}
			wholesaleRaw, err := strconv.Atoi(dealerCell)
			if err != nil {
				return fmt.Errorf("невалидная дилерская цена %q: %w", dealerCell, err)
			}

			if err := repository.UpdateGateSizePrice(
				width, height, gateType,
				decimal.NewFromInt(int64(wholesaleRaw)),
				decimal.NewFromInt(int64(retailRaw)),
			); err != nil {
				return fmt.Errorf("ошибка обновления размера %dx%d: %w", width, height, err)
			}
		}
	}

	return nil
}
