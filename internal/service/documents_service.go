package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/samborkent/uuidv7"

	"github.com/Euhcslel/SagaWeb/internal/database"
	"github.com/Euhcslel/SagaWeb/internal/docgen"
	"github.com/Euhcslel/SagaWeb/internal/domain/enums"
	"github.com/Euhcslel/SagaWeb/internal/domain/industrial_gate_drives"
	"github.com/Euhcslel/SagaWeb/internal/domain/options"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_gates"
	"github.com/Euhcslel/SagaWeb/internal/domain/order_products"
	"github.com/Euhcslel/SagaWeb/internal/domain/users"
	errs "github.com/Euhcslel/SagaWeb/internal/errors"
	"github.com/Euhcslel/SagaWeb/internal/generated"
	"github.com/Euhcslel/SagaWeb/internal/repository"
	"github.com/Euhcslel/SagaWeb/internal/types"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func UploadOfferToOrder(user *users.User, orderID int64, file multipart.File, handler *multipart.FileHeader) error {
	return uploadDocument(user, orderID, file, handler, uploadConfig{
		dirEnvKey: "OFFERS_DIRECTORY",
		attach:    repository.AttachOfferToOrder,
	})
}

func UploadContractToOrder(user *users.User, orderID int64, file multipart.File, handler *multipart.FileHeader) error {
	return uploadDocument(user, orderID, file, handler, uploadConfig{
		dirEnvKey: "CONTRACTS_DIRECTORY",
		attach:    repository.AttachAppendiceToOrder,
	})
}

func UploadBillToOrder(user *users.User, orderID int64, file multipart.File, handler *multipart.FileHeader) error {
	return uploadDocument(user, orderID, file, handler, uploadConfig{
		dirEnvKey: "BILLS_DIRECTORY",
		attach:    repository.AttachBillToOrder,
	})
}

func UploadDocumentToOrder(user *users.User, orderID int64, file multipart.File, handler *multipart.FileHeader) error {
	return uploadDocument(user, orderID, file, handler, uploadConfig{
		dirEnvKey: "DOCUMENTS_DIRECTORY",
		attach:    repository.AttachDocumentToOrder,
	})
}

// Структура конфигурации загружаемого файла
// Включает в себя путь до каталога с файлом и функцию репозитория, которая прикрепляет документ к заказу
type uploadConfig struct {
	dirEnvKey string
	attach    func(db *gorm.DB, orderID int64, name, filename string) error
}

func uploadDocument(user *users.User, orderID int64,
	file multipart.File, handler *multipart.FileHeader, cfg uploadConfig) error {

	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}

	dir := os.Getenv(cfg.dirEnvKey)
	if dir == "" {
		return fmt.Errorf("переменная %s не задана", cfg.dirEnvKey)
	}

	ext := filepath.Ext(handler.Filename)
	name := strings.TrimSuffix(handler.Filename, ext)
	storedFilename := uuidv7.New().String() + "_" + filepath.Base(handler.Filename)

	path := filepath.Join(dir, storedFilename)
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		os.Remove(path)
		return err
	}

	if err = cfg.attach(database.DB, orderID, name, storedFilename); err != nil {
		os.Remove(path)
		return err
	}

	return nil
}

func GetAllOrderDocuments(user *users.User, orderID int64) (*generated.DocumentsList, error) {
	if err := checkOrderAccess(user, orderID); err != nil {
		return nil, err
	}

	var resp *generated.DocumentsList

	offersNumberList, err := repository.GetOffersNumberList(database.DB, orderID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	offers := make([]*generated.Offer, 0, len(offersNumberList))
	for _, number := range offersNumberList {
		offers = append(offers, &generated.Offer{
			OfferNumber: number,
		})
	}

	contractsNumberList, err := repository.GetAppendicesNumberList(database.DB, orderID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	contracts := make([]*generated.Contract, 0, len(contractsNumberList))
	for _, number := range contractsNumberList {
		contracts = append(contracts, &generated.Contract{
			ContractNumber: number,
		})
	}

	billsNumberList, err := repository.GetBillsNumberList(database.DB, orderID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	bills := make([]*generated.Bill, 0, len(billsNumberList))
	for _, number := range billsNumberList {
		bills = append(bills, &generated.Bill{
			BillNumber: number,
		})
	}

	documentsNameList, err := repository.GetDocumentsNameList(database.DB, orderID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	documents := make([]*generated.Document, 0, len(documentsNameList))
	for _, name := range documentsNameList {
		documents = append(documents, &generated.Document{
			Name: name,
		})
	}

	resp = &generated.DocumentsList{
		Documents: documents,
		Offers:    offers,
		Bills:     bills,
		Contracts: contracts,
	}

	return resp, nil
}

// Структура, описывающая файл
// Включает имя и путь до файла
type FileInfo struct {
	FilePath string
	FileName string
}

func GetFileInfo(user *users.User, orderID int64, docType string, docName string) (*FileInfo, error) {
	err := checkOrderAccess(user, orderID)
	if err != nil {
		return nil, err
	}

	fileInfo := FileInfo{}
	fileInfo.FilePath, fileInfo.FileName, err = resolveDocumentPath(orderID, docType, docName)
	if err != nil {
		return nil, err
	}

	return &fileInfo, nil
}

func DeleteOrderDocument(user *users.User, orderID int64, docType string, docName string) error {
	if err := checkOrderAccess(user, orderID); err != nil {
		return err
	}

	filePath, _, err := resolveDocumentPath(orderID, docType, docName)
	if err != nil {
		return err
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	switch enums.DocumentType(docType) {
	case enums.BillDocumentType:
		if err := repository.DeleteOrderBill(tx, orderID, docName); err != nil {
			return err
		}
	case enums.ContractDocumentType:
		if err := repository.DeleteOrderAppendice(tx, orderID, docName); err != nil {
			return err
		}
	case enums.OfferDocumentType:
		if err := repository.DeleteOrderOffer(tx, orderID, docName); err != nil {
			return err
		}
	case enums.OtherDocumentType:
		if err := repository.DeleteOrderDocument(tx, orderID, docName); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	os.Remove(filePath)
	return nil
}

func resolveDocumentPath(orderID int64, docType string, docName string) (filePath string, fileName string, err error) {
	var dir string

	switch enums.DocumentType(docType) {
	case enums.BillDocumentType:
		fileName, err = repository.GetBillFileName(orderID, docName)
		dir = os.Getenv("BILLS_DIRECTORY")

	case enums.OfferDocumentType:
		fileName, err = repository.GetOfferFileName(orderID, docName)
		dir = os.Getenv("OFFERS_DIRECTORY")

	case enums.ContractDocumentType:
		fileName, err = repository.GetAppendiceFileName(orderID, docName)
		dir = os.Getenv("CONTRACTS_DIRECTORY")

	case enums.OtherDocumentType:
		fileName, err = repository.GetDocumentFileName(orderID, docName)
		dir = os.Getenv("DOCUMENTS_DIRECTORY")

	default:
		return "", "", errs.ErrInvalidDocumentType
	}

	if err != nil {
		return "", "", err
	}

	if fileName == "" {
		return "", "", errs.ErrNotFound
	}

	if dir == "" {
		return "", "", fmt.Errorf("directory env for document type %s is not set", docType)
	}

	return filepath.Join(dir, fileName), fileName, nil
}

func GetOfferForOrder(
	user *users.User,
	orderData *generated.OrderRequest,
) ([]byte, error) {
	order := types.Order{}

	for i, gate := range orderData.OrderGates {
		gateType, err := enums.GateTypeFromProto(gate.GateType)
		if err != nil {
			return nil, err
		}

		driveType, err := enums.GetDriveTypeFromProto(gate.Drive)
		if err != nil {
			return nil, err
		}

		// справочники: цвет, тип подъёма, кол-во циклов
		color, err := repository.GetColorByID(database.DB, gate.ColorOutId)
		if err != nil {
			return nil, err
		}

		liftType, err := repository.GetLiftTypeByID(database.DB, gate.LiftTypeId)
		if err != nil {
			return nil, err
		}

		cycleAmount, err := repository.GetCycleAmountByID(database.DB, gate.CycleAmountId)
		if err != nil {
			return nil, err
		}

		// привод — три варианта
		var drive any
		switch driveType {
		case enums.IndDriveType:
			d, err := repository.GetIndustrialDriveByID(database.DB, gate.Drive.GetIndustrial().DriveId)
			if err != nil {
				return nil, err
			}
			drive = *d
		case enums.ResDriveType:
			d, err := repository.GetResidentialDriveByID(database.DB, gate.Drive.GetResidential().DriveId)
			if err != nil {
				return nil, err
			}
			rail, err := repository.GetRailByID(database.DB, gate.Drive.GetResidential().RailId)
			if err != nil {
				return nil, err
			}
			drive = types.ResidentialDriveRail{
				Drive: *d,
				Rail:  *rail,
			}
		case enums.ManualDriveType:
			info, err := repository.GetManualDrivePrice(database.DB)
			if err != nil {
				return nil, err
			}
			drive = types.ManualDrive{
				ChainLength: gate.Drive.GetManual().ChainLength,
				PriceInfo:   *info,
			}
		}

		// опции: грузим объекты + считаем сумму
		gateOptions, optionsPrices, err := loadOptions(gate.Options)
		if err != nil {
			return nil, err
		}

		// цены привода для общего расчёта ворот
		drivePrices := drivePricePair(drive)

		gatePrices, err := calculateGatePrice(
			int64(gate.Width), int64(gate.Height), gateType,
			drivePrices, *liftType, *cycleAmount, optionsPrices,
		)
		if err != nil {
			return nil, err
		}

		// собираем OrderGate со связями (для печати имён)
		gateConfiguration := order_gates.OrderGate{
			RowNumber:          int64(i + 1),
			GateType:           gateType,
			Width:              gate.Width,
			Height:             gate.Height,
			Headroom:           gate.Headroom,
			LiftTypeID:         gate.LiftTypeId,
			ColorOutID:         gate.ColorOutId,
			CycleAmountID:      gate.CycleAmountId,
			DriveType:          driveType,
			Amount:             gate.Amount,
			GateRetailPrice:    gatePrices.RetailPrice,
			GateWholesalePrice: gatePrices.WholesalePrice,
			LiftType:           *liftType,
			ColorOut:           *color,
			CycleAmount:        *cycleAmount,
		}

		order.Gates = append(order.Gates, types.Gate{
			Gate:    gateConfiguration,
			Drive:   drive,
			Options: gateOptions,
		})
	}

	// --- товары ---
	products, err := loadProducts(orderData.Products)
	if err != nil {
		return nil, err
	}
	order.Products = products

	return docgen.GenerateClientOffer(&order)
}

// buildOrderForDoc загружает ворота, приводы, опции и товары заказа для генерации документов.
func buildOrderForDoc(orderID int64) (types.Order, error) {
	orderGates, err := repository.GetOrderGatesByOrderID(database.DB, orderID)
	if err != nil {
		return types.Order{}, err
	}

	var gateIDs []int64
	for _, g := range orderGates {
		gateIDs = append(gateIDs, g.RowNumber)
	}

	gateOptionsRows, err := repository.GetAllOptionsForOrder(database.DB, orderID, gateIDs)
	if err != nil {
		return types.Order{}, err
	}
	optionsMap := make(map[int64][]options.Option)
	for _, opt := range gateOptionsRows {
		optionsMap[opt.RowNumber] = append(optionsMap[opt.RowNumber], opt.Option)
	}

	var order types.Order
	for _, og := range orderGates {
		drive, err := loadDriveForGate(og)
		if err != nil {
			return types.Order{}, err
		}
		order.Gates = append(order.Gates, types.Gate{
			Gate:    og,
			Drive:   drive,
			Options: optionsMap[og.RowNumber],
		})
	}

	order.Products, err = repository.GetAllOrderProducts(database.DB, orderID)
	if err != nil {
		return types.Order{}, err
	}

	orderStatus, err := repository.GetOrderStatus(orderID)
	if err != nil {
		return types.Order{}, err
	}
	if enums.OrderStatus(*orderStatus) != enums.OrderStatusPending {
		finalizedAt, err := repository.GetOrderFinalizedAt(database.DB, orderID)
		if err != nil {
			return types.Order{}, err
		}
		if finalizedAt != nil && len(order.Products) > 0 {
			productIDs := make([]int64, len(order.Products))
			for i, p := range order.Products {
				productIDs[i] = p.ProductID
			}
			historicalPrices, err := repository.GetProductHistoricalPrices(database.DB, productIDs, *finalizedAt)
			if err != nil {
				return types.Order{}, err
			}
			for i, p := range order.Products {
				if prices, ok := historicalPrices[p.ProductID]; ok {
					order.Products[i].Product.RetailPrice = prices.RetailPrice
					order.Products[i].Product.WholesalePrice = prices.WholesalePrice
				}
			}
		}
	}

	return order, nil
}

// GetDealerOfferForClient генерирует КП дилера для покупателя (только розничные цены).
func GetDealerOfferForClient(user *users.User, orderID int64, clientName string) ([]byte, error) {
	if user.Role != enums.DealerRole {
		return nil, errs.ErrForbidden
	}
	if err := checkOrderAccess(user, orderID); err != nil {
		return nil, err
	}

	order, err := buildOrderForDoc(orderID)
	if err != nil {
		return nil, err
	}

	dealer, err := repository.GetDealerInfo(database.DB, user.ID)
	if err != nil {
		return nil, err
	}
	supplier := docgen.SupplierInfo{
		Name:  dealer.Company.Name,
		Phone: user.PhoneNumber,
		Email: user.Email,
	}

	return docgen.GenerateDealerOfferForClient(&order, &orderID, supplier, clientName)
}

// GetDealerOfferForSelf генерирует КП дилера для себя (розничные + дилерские цены).
func GetDealerOfferForSelf(user *users.User, orderID int64, clientName string) ([]byte, error) {
	if user.Role != enums.DealerRole {
		return nil, errs.ErrForbidden
	}
	if err := checkOrderAccess(user, orderID); err != nil {
		return nil, err
	}

	order, err := buildOrderForDoc(orderID)
	if err != nil {
		return nil, err
	}

	dealer, err := repository.GetDealerInfo(database.DB, user.ID)
	if err != nil {
		return nil, err
	}
	supplier := docgen.SupplierInfo{
		Name:  dealer.Company.Name,
		Phone: user.PhoneNumber,
		Email: user.Email,
	}

	return docgen.GenerateDealerOfferForSelf(&order, &orderID, supplier, clientName)
}

// loadOptions грузит опции по id и возвращает их + сумму.
func loadOptions(gateOptions []*generated.Option) ([]options.Option, types.PricePair, error) {
	if len(gateOptions) == 0 {
		return nil, types.PricePair{}, nil
	}

	ids := make([]int64, 0, len(gateOptions))
	for _, o := range gateOptions {
		if o.OptionId == 0 || o.Amount <= 0 {
			continue
		}
		ids = append(ids, o.OptionId)
	}
	if len(ids) == 0 {
		return nil, types.PricePair{}, nil
	}

	list, err := repository.GetOptionsByIDs(database.DB, ids)
	if err != nil {
		return nil, types.PricePair{}, err
	}

	byID := make(map[int64]options.Option, len(list))
	for _, o := range list {
		byID[o.ID] = o
	}

	var sum types.PricePair
	result := make([]options.Option, 0, len(gateOptions))
	for _, o := range gateOptions {
		if o.OptionId == 0 || o.Amount <= 0 {
			continue
		}
		opt, ok := byID[o.OptionId]
		if !ok {
			return nil, types.PricePair{}, errors.New("option not found")
		}
		amount := decimal.NewFromInt32(o.Amount)
		sum.RetailPrice = sum.RetailPrice.Add(opt.RetailPrice.Mul(amount))
		sum.WholesalePrice = sum.WholesalePrice.Add(opt.WholesalePrice.Mul(amount))
		result = append(result, opt)
	}
	return result, sum, nil
}

// loadProducts грузит товары и возвращает OrderProduct с заполненным Product.
func loadProducts(in []*generated.Product) ([]order_products.OrderProduct, error) {
	result := make([]order_products.OrderProduct, 0, len(in))
	for _, p := range in {
		if p.ProductId == 0 || p.Amount <= 0 {
			continue
		}
		prod, err := repository.GetProductByID(database.DB, p.ProductId)
		if err != nil {
			return nil, err
		}
		result = append(result, order_products.OrderProduct{
			ProductID: p.ProductId,
			Amount:    p.Amount,
			Product:   *prod,
		})
	}
	return result, nil
}

// GetAppendiceForOrder загружает заказ из БД и генерирует PDF приложения к договору.
// Номер договора и дата берутся из записи dealer_contract; номер приложения — count(существующих) + 1.
func GetAppendiceForOrder(user *users.User, orderID int64) ([]byte, int, error) {
	if err := checkOrderAccess(user, orderID); err != nil {
		return nil, 0, err
	}

	orderRecord, err := repository.GetOrderWithCustomer(database.DB, orderID)
	if err != nil {
		return nil, 0, err
	}

	if orderRecord.DealerID == nil {
		return nil, 0, errs.ErrNoDealerContract
	}

	contract, err := repository.GetDealerContract(database.DB, *orderRecord.DealerID)
	if err != nil {
		return nil, 0, err
	}
	if contract == nil {
		return nil, 0, errs.ErrNoDealerContract
	}

	existing, err := repository.GetAppendicesNumberList(database.DB, orderID)
	if err != nil {
		return nil, 0, err
	}
	appendiceNumber := len(existing) + 1

	order, err := buildOrderForDoc(orderID)
	if err != nil {
		return nil, 0, err
	}

	var customer docgen.CustomerInfo
	if orderRecord.Dealer != nil {
		customer.OrganizationName = orderRecord.Dealer.Company.Name
		customer.ContactPerson = orderRecord.Dealer.User.Fullname
		customer.Phone = orderRecord.Dealer.User.PhoneNumber
	} else {
		customer.ContactPerson = orderRecord.Manager.Fullname
		customer.Phone = orderRecord.Manager.PhoneNumber
	}

	pdf, err := docgen.GenerateAppendice(&order, customer, appendiceNumber, contract.ContractNumber, contract.SignedAt)
	return pdf, appendiceNumber, err
}

// loadDriveForGate загружает привод для конкретных ворот в зависимости от типа
func loadDriveForGate(og order_gates.OrderGate) (any, error) {
	switch og.DriveType {
	case enums.IndDriveType:
		d, err := repository.GetIndustrialDriveForGate(database.DB, og.OrderID, og.RowNumber)
		if err != nil {
			return nil, err
		}
		ind, err := repository.GetIndustrialDriveByID(database.DB, d.DriveID)
		if err != nil {
			return nil, err
		}
		return *ind, nil
	case enums.ResDriveType:
		rd, err := repository.GetResidentialDriveForGate(database.DB, og.OrderID, og.RowNumber)
		if err != nil {
			return nil, err
		}
		drv, err := repository.GetResidentialDriveByID(database.DB, rd.DriveID)
		if err != nil {
			return nil, err
		}
		rail, err := repository.GetRailByID(database.DB, rd.RailID)
		if err != nil {
			return nil, err
		}
		return types.ResidentialDriveRail{Drive: *drv, Rail: *rail}, nil
	case enums.ManualDriveType:
		md, err := repository.GetManualDriveForGate(database.DB, og.OrderID, og.RowNumber)
		if err != nil {
			return nil, err
		}
		info, err := repository.GetManualDrivePrice(database.DB)
		if err != nil {
			return nil, err
		}
		return types.ManualDrive{ChainLength: md.ChainLength, PriceInfo: *info}, nil
	}
	return nil, nil
}

// drivePricePair достаёт пару цен из любого варианта привода.
func drivePricePair(drive any) types.PricePair {
	switch d := drive.(type) {
	case industrial_gate_drives.IndustrialGateDrive:
		return types.PricePair{RetailPrice: d.RetailPrice, WholesalePrice: d.WholesalePrice}
	case types.ResidentialDriveRail:
		return types.PricePair{
			RetailPrice:    d.Drive.RetailPrice.Add(d.Rail.RetailPrice),
			WholesalePrice: d.Drive.WholesalePrice.Add(d.Rail.WholesalePrice),
		}
	case types.ManualDrive:
		chain := decimal.NewFromInt32(d.ChainLength)
		return types.PricePair{
			RetailPrice:    d.PriceInfo.ChainMeterRetailPrice.Mul(chain).Add(d.PriceInfo.RcpRetailPrice),
			WholesalePrice: d.PriceInfo.ChainMeterWholesalePrice.Mul(chain).Add(d.PriceInfo.RcpWholesalePrice),
		}
	}
	return types.PricePair{}
}
