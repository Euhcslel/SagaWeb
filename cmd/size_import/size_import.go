package main

import (
	"fmt"
	"log"
	"project/pkg/database"
	"project/pkg/models"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
)

func main() {
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatal("Error loading .env file")
	}
	db := database.InitDB()
	var homeId int
	var promId int

	db.Model(&models.GateType{}).
		Where("name = ?", "Бытовые ворота").
		Select("id").
		Scan(&homeId)

	db.Model(&models.GateType{}).
		Where("name = ?", "Промышленные ворота").
		Select("id").
		Scan(&promId)

	insertHomeSizes(homeId)
	insertPromSizes(promId)
}

func insertPromSizes(id int) {
	createSizes("C:/Users/np_gl/Desktop/Диплом/sizes_import/prom_dealer_sizes.xlsx", "C:/Users/np_gl/Desktop/Диплом/sizes_import/prom_client_sizes.xlsx", id)
}

func insertHomeSizes(id int) {
	createSizes("C:/Users/np_gl/Desktop/Диплом/sizes_import/home_dealer_sizes.xlsx", "C:/Users/np_gl/Desktop/Диплом/sizes_import/home_client_sizes.xlsx", id)
}

func createSizes(dealerPath string, clientPath string, gateType int) {
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

	db := database.InitDB()
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

			size.GateTypeID = int64(gateType)

			sizes = append(sizes, size)
		}
	}

	db.Create(&sizes)
}

func readExcelFile(path string) *excelize.File {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return f
}
