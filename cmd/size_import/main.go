package main

import (
	"fmt"
	"log"
	"os"
	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/gate_sizes"

	"strconv"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

func main() {
	err_load := godotenv.Load("../../.env")
	if err_load != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	dealerPathInd := os.Getenv("EXCEL_IMPORT_IND_DEALER_PATH")
	if dealerPathInd == "" {
		log.Fatal("EXCEL_IMPORT_IND_DEALER_PATH не найден в .env файле")
	}

	clientPathInd := os.Getenv("EXCEL_IMPORT_IND_CLIENT_PATH")
	if clientPathInd == "" {
		log.Fatal("EXCEL_IMPORT_IND_CLIENT_PATH не найден в .env файле")
	}

	dealerPathRes := os.Getenv("EXCEL_IMPORT_RES_DEALER_PATH")
	if dealerPathRes == "" {
		log.Fatal("EXCEL_IMPORT_RES_DEALER_PATH не найден в .env файле")
	}

	clientPathRes := os.Getenv("EXCEL_IMPORT_RES_CLIENT_PATH")
	if clientPathRes == "" {
		log.Fatal("EXCEL_IMPORT_RES_CLIENT_PATH не найден в .env файле")
	}

	if err := createSizes(dealerPathInd, clientPathInd, enums.GateTypeInd); err != nil {
		log.Printf("Ошибка при попытке добавить цены для промышленных ворот: %v", err)
	}

	if err := createSizes(dealerPathRes, clientPathRes, enums.GateTypeRes); err != nil {
		log.Printf("Ошибка при попытке добавить цены для бытовых ворот: %v", err)
	}
}

func createSizes(dealerPath string, clientPath string, gateType enums.GateType) (err error) {
	clientPricesFile, err := excelize.OpenFile(clientPath)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении файла с клиентскими ценами: %w", err)
	}
	defer func() {
		if chain_err := clientPricesFile.Close(); chain_err != nil {
			err = fmt.Errorf("file close error: %w\nwith error:%w", chain_err, err)
		}
	}()

	dealerPricesFile, err := excelize.OpenFile(dealerPath)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении файла с дилерскими ценами: %w", err)
	}
	defer func() {
		if chain_err := dealerPricesFile.Close(); chain_err != nil {
			err = fmt.Errorf("file close error: %w\nwith error:%w", chain_err, err)
		}
	}()

	sizeList := []gate_sizes.GateSize{}

	clientRows, err := clientPricesFile.GetRows("Лист1")
	if err != nil {
		return fmt.Errorf("Ошибка при попытке прочтать строки в файле с клиентскими ценами: %w", err)
	}

	clientCols, err := clientPricesFile.GetCols("Лист1")
	if err != nil {
		return fmt.Errorf("Ошибка при попытке прочитать столбцы в файле с клиентскими ценами: %w", err)
	}

	dealerRows, err := dealerPricesFile.GetRows("Лист1")
	if err != nil {
		return fmt.Errorf("Ошибка при попытке прочитать строки в файле с дилерскими ценами: %w", err)
	}

	dealerCols, err := dealerPricesFile.GetCols("Лист1")
	if err != nil {
		return fmt.Errorf("Ошибка при попытке прочитать столбцы в файле с дилерскими ценами: %w", err)
	}

	if len(dealerCols) != len(clientCols) || len(dealerRows) != len(clientRows) {
		return fmt.Errorf("Число столбцов или строк в файлах различаются!")
	}

	for i := 1; i < len(clientRows); i++ {
		size := gate_sizes.GateSize{}

		height, err := clientPricesFile.GetCellValue("Лист1", "A"+strconv.Itoa(i))
		if height == "" {
			continue
		} else if err != nil {
			return fmt.Errorf("Ошибка при попытке прочитать значение ячейки высоты: %w", err)
		}
		intHeight, err := strconv.Atoi(height)
		if err != nil {
			return fmt.Errorf("Ошибка при попытке преобразовать значение высоты: %w", err)
		}
		size.Height = int64(intHeight)

		for j := 1; j < len(clientCols); j++ {
			if j == 1 || i == 1 {
				continue
			}
			columnLetter, err := excelize.ColumnNumberToName(j)
			if err != nil {
				return fmt.Errorf("Ошибка при попытке преобразовать значение к букве в Excel: %w", err)
			}

			width, err := clientPricesFile.GetCellValue("Лист1", columnLetter+"1")
			if width == "" {
				continue
			} else if err != nil {
				return fmt.Errorf("Ошибка при попытке прочитать значение ячейки ширины: %w", err)
			}
			intWidth, err := strconv.Atoi(width)
			if err != nil {
				return fmt.Errorf("Ошибка при попытке преобразовать значение ширины: %w", err)
			}
			size.Width = int64(intWidth)

			clientCell, err := clientPricesFile.GetCellValue("Лист1", columnLetter+strconv.Itoa(i))
			if err != nil {
				return fmt.Errorf("Ошибка при попытке прочитать значение клиентской цены: %w", err)
			}

			dealerCell, err := dealerPricesFile.GetCellValue("Лист1", columnLetter+strconv.Itoa(i))
			if err != nil {
				return fmt.Errorf("Ошибка при попытке прочитать значение дилерской цены: %w", err)
			}

			if clientCell == "" || dealerCell == "" {
				log.Printf("Значение для ширины %s и высоты %s отсутствует \n", width, height)
				continue
			}

			intWholesalePrice, err := strconv.Atoi(dealerCell)
			if err != nil {
				return fmt.Errorf("Ошибка при попытке преобразовать значение дилерской цены: %w", err)
			}
			size.WholesalePrice = decimal.NewFromInt(int64(intWholesalePrice))

			intRetailPrice, err := strconv.Atoi(clientCell)
			if err != nil {
				return fmt.Errorf("Ошибка при попытке преобразовать значение клиентской цены: %w", err)
			}
			size.RetailPrice = decimal.NewFromInt(int64(intRetailPrice))

			size.GateType = gateType

			sizeList = append(sizeList, size)
		}
	}

	if err := database.DB.Create(&sizeList).Error; err != nil {
		return fmt.Errorf("Ошибка при попытке создать значения цен в базе данных: %w", err)
	}

	return nil
}
