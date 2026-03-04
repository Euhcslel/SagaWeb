package main

import (
	"fmt"
	"log"
	"os"
	"project/pkg/database"
	"project/pkg/models"
	"strconv"

	"github.com/joho/godotenv"
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

	createSizes(dealerPathInd, clientPathInd, models.GateTypeInd)
	createSizes(dealerPathRes, clientPathRes, models.GateTypeRes)
}

func createSizes(dealerPath string, clientPath string, gateType models.GateType) {
	clientPricesFile := readExcelFile(clientPath)
	if clientPricesFile == nil {
		return
	}
	defer func() {
		if err := clientPricesFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	dealerPricesFile := readExcelFile(dealerPath)
	if dealerPricesFile == nil {
		return
	}
	defer func() {
		if err := dealerPricesFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sizes := []models.Size{}

	clientRows, err := clientPricesFile.GetRows("Лист1")
	if err != nil {
		fmt.Println(err)
		return
	}

	clientCols, err := clientPricesFile.GetCols("Лист1")
	if err != nil {
		fmt.Println(err)
		return
	}

	dealerRows, err := dealerPricesFile.GetRows("Лист1")
	if err != nil {
		fmt.Println(err)
		return
	}

	dealerCols, err := dealerPricesFile.GetCols("Лист1")
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(dealerCols) != len(clientCols) || len(dealerRows) != len(clientRows) {
		fmt.Println("Число столбцов или строк в файлах Prom различаются!")
		return
	}

	for i := 1; i < len(clientRows); i++ {
		size := models.Size{}

		height, _ := clientPricesFile.GetCellValue("Лист1", "A"+strconv.Itoa(i))
		intHeight, _ := strconv.Atoi(height)
		size.Height = int64(intHeight)

		for j := 1; j < len(clientCols); j++ {
			if j == 1 || i == 1 {
				continue
			}
			columnLetter, err := excelize.ColumnNumberToName(j)
			if err != nil {
				fmt.Println(err)
				return
			}

			width, _ := clientPricesFile.GetCellValue("Лист1", columnLetter+"1")
			intWidth, _ := strconv.Atoi(width)
			size.Width = int64(intWidth)

			clientCell, err := clientPricesFile.GetCellValue("Лист1", columnLetter+strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
				return
			}

			dealerCell, err := dealerPricesFile.GetCellValue("Лист1", columnLetter+strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
				return
			}

			if clientCell == "" || dealerCell == "" {
				log.Printf("Значение для ширины %s и высоты %s отсутствует \n", width, height)
				continue
			}

			intWholesalePrice, _ := strconv.Atoi(dealerCell)
			size.WholesalePrice = int64(intWholesalePrice)

			intRetailPrice, _ := strconv.Atoi(clientCell)
			size.RetailPrice = int64(intRetailPrice)

			size.GateType = gateType

			sizes = append(sizes, size)
		}
	}

	database.DB.Create(&sizes)
}

func readExcelFile(path string) *excelize.File {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return f
}
