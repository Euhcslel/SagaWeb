package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"project/pkg/database"
	"project/pkg/helpers"
	"project/pkg/models"
	"strconv"
	"sync"

	"project/pkg/types"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var cssPath = "/assets/styles/"

// Route: /log
// Method: GET
func GetSignInForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"css": cssPath + "auth.css",
	}

	if err := templates.ExecuteTemplate(w, "log.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /log
// Method: POST
func SignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user models.User
	err = database.DB.Model(&models.User{}).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}
	userId := user.ID
	dbPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	helpers.SetSession(w, userId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /reg
// Method: GET
func GetSignUpFrom(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"css": cssPath + "auth.css",
	}
	if err := templates.ExecuteTemplate(w, "reg.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /reg
// Method: POST
func SignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusBadRequest)
		return
	}

	accountType := r.FormValue("accountType")

	switch accountType {
	case "dealer":
		fullname := r.FormValue("fullname")
		company := r.FormValue("company")
		phone := r.FormValue("phone")
		email := r.FormValue("email")

		var status models.Status
		if err := database.DB.Model(&models.Status{}).Where("name = ?", "Ожидает подтверждения").Find(&status).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}

		request := models.DealerRegRequest{
			Company:     company,
			Fullname:    fullname,
			PhoneNumber: phone,
			Email:       email,
			StatusID:    int64(status.ID),
		}
		database.DB.Create(&request)

	case "client":
		username := r.FormValue("username")
		fullname := r.FormValue("fullname")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		password := r.FormValue("password")

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
			return
		}

		var role models.Role
		if err := database.DB.Model(&models.Role{}).Where("name = ?", "client").Find(&role).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}

		user := models.User{
			Fullname:    fullname,
			PhoneNumber: phone,
			Email:       email,
			Password:    string(passwordHash),
			Username:    username,
			RoleID:      role.ID,
		}
		database.DB.Create(&user)

		helpers.SetSession(w, user.ID)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Route: /
// Method: GET
func MainHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		user = models.User{}
	} else {
		token := sessionToken.Value
		user = helpers.GetUserBySessionToken(token)
	}

	data := map[string]any{
		"css":  "/assets/styles/main.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "main.html", data); err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
	}
}

// Route: /contacts
// Method: GET
func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "contacts.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route: /user
// Method: GET
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var user models.User
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	} else {
		token := sessionToken.Value
		user = helpers.GetUserBySessionToken(token)
	}

	data := map[string]any{
		"css":  cssPath + "user.css",
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "user.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
	}
}

// Route: /user/dealers
// Method: GET
func GetUserDealers(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)
	role := user.Role.Name

	if role == "manager" {
		var dealers []models.ManagerAndDealer
		if err := database.DB.Model(models.ManagerAndDealer{}).Preload("Dealer").Where("manager_id = ?", user.ID).Find(&dealers).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"css":     "",
			"user":    user,
			"dealers": dealers,
		}

		if err := templates.ExecuteTemplate(w, "dealers.html", data); err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
	}

}

// Route: /orders
// Method: GET
func GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)

	var orders []models.Sale
	err = database.DB.
		Preload("Client").
		Preload("Manager").
		Where("client_id = ?", user.ID).
		Limit(5).
		Find(&orders).Error
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":    "",
		"orders": orders,
	}

	if err := templates.ExecuteTemplate(w, "orders.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /api/orders
// Method: GET
func GetAllUserOrdersAPI(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: GET
func GetUserOrderById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderId := vars["order_id"]

	// Получаем все ворота выбранного заказа
	var orderGates []models.SalesAndGate
	database.DB.Preload("GateType").
		Preload("LiftType").
		Preload("ColorIn").
		Preload("ColorOut").
		Preload("CycleAmount").
		Preload("Status").
		Model(models.SalesAndGate{}).
		Where("sale_id = ?", orderId).
		Find(&orderGates)

	// Собираем все ID ворот
	var gateIDs []int64
	for _, gate := range orderGates {
		gateIDs = append(gateIDs, gate.RowNumber)
	}

	// Получаем все опции всех ворот
	var options []models.GatesAndSalesOption
	database.DB.
		Where("row_number IN ?", gateIDs).
		Preload("Option").
		Find(&options)

	// Группируем опции по ID ворот
	optionsMap := make(map[int64][]models.Option)
	for _, opt := range options {
		optionsMap[opt.RowNumber] = append(optionsMap[opt.RowNumber], opt.Option)
	}

	// Группируем в gates
	order := types.Order{
		Gates:    []types.Gate{},
		Products: []models.Product{},
	}
	for _, og := range orderGates {
		order.Gates = append(order.Gates, types.Gate{
			Gate:    og,
			Options: optionsMap[og.RowNumber],
		})
	}

	data := map[string]any{
		"css":    "",
		"orders": order,
	}

	if err := templates.ExecuteTemplate(w, "order.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}
}

// Route: /orders/{order_id}/{gate_id}
// Method: GET
func GetGateInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/{gate_id}/options
// Method: GET
func GetGateOptions(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/documents
// Method: GET
func GetOrderDocuments(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/products
// Method: GET
func GetProductsInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: DELETE
func DeleteUserOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}
// Method: POST
func AddNewGateInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders
// Method: POST
func CreateNewOrder(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)
	role := user.Role.Name

	var order models.Sale
	if role == "client" || role == "dealer" {
		var managerID int64
		err = database.DB.Model(&models.ManagerAndDealer{}).Select("manager_id").Where("dealer_id = ?", user.ID).First(&managerID).Error
		if err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}

		order = models.Sale{
			ClientID:  user.ID,
			ManagerID: managerID,
		}
		database.DB.Create(&order)
	} else {
		order = models.Sale{
			ManagerID: user.ID,
		}
		database.DB.Create(&order)
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusBadRequest)
		return
	}

	orderGatesJSON := r.FormValue("orderGates")
	productsJSON := r.FormValue("products")

	var orderRequest types.OrderRequest

	if err := json.Unmarshal([]byte(orderGatesJSON), &orderRequest.OrderGates); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse orderGates: %v", err), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(productsJSON), &orderRequest.Products); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse products: %v", err), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	for i, gate := range orderRequest.OrderGates {
		gateTypeInt, _ := strconv.ParseInt(gate.GateType, 10, 32)
		width := gate.Width
		height := gate.Height
		liftTypeID, _ := strconv.ParseInt(gate.LiftType, 10, 64)
		colorInID, _ := strconv.ParseInt(gate.ColorIn, 10, 64)
		colorOutID, _ := strconv.ParseInt(gate.ColorOut, 10, 64)
		cycleAmountID, _ := strconv.ParseInt(gate.CycleAmount, 10, 64)

		orderDetails := models.SalesAndGate{
			SaleID:        order.ID,
			RowNumber:     int64(i + 1),
			GateTypeID:    int32(gateTypeInt),
			Width:         int32(width),
			Height:        int32(height),
			LiftTypeID:    liftTypeID,
			ColorInID:     colorInID,
			ColorOutID:    colorOutID,
			CycleAmountID: cycleAmountID,
			TotalPrice:    int32(gate.GatePrice),
			StatusID:      9,
		}

		wg.Go(func() {
			database.DB.Create(&orderDetails)
		})
	}

	wg.Wait()
	http.Redirect(w, r, "/orders", http.StatusSeeOther)
}

// Route: /orders/{order_id}/products
// Method: POST
func AddNewProductInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/products/{product_id}
// Method: PUT
func UpdateProductList(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/products/{product_id}
// Method: DELETE
func DeleteProductFromOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/{gate_id}
// Method: DELETE
func DeleteGateFromOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /orders/{order_id}/{gate_id}
// Method: PUT
func UpdateGateInOrder(w http.ResponseWriter, r *http.Request) {

}

// Route: /gate_types
// Method: GET
func GetGateTypesList(w http.ResponseWriter, r *http.Request) {
	var gateTypes []models.GateType
	err := database.DB.Find(&gateTypes).Error
	if err != nil {
		helpers.WriteErrorDebug(w, err, http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"css":       "",
		"gateTypes": gateTypes,
	}

	if err := templates.ExecuteTemplate(w, "gate_type.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
	}
}

// Route: /calculator
// Method: GET
func GetCalculatorForUser(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	var user models.User
	if err == nil {
		token := sessionToken.Value
		user = helpers.GetUserBySessionToken(token)
	}

	query := r.URL.Query()
	gateType := query.Get("gateType")

	cfg := types.Config{
		LiftTypes:    []models.LiftType{},
		Colors:       []models.Color{},
		WidthParams:  types.WidthParams{},
		HeightParams: types.HeightParams{},
	}

	var gateTypes []models.GateType
	if err = database.DB.Find(&gateTypes).Error; err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
		return
	}

	gateTypesMap := make(map[int32]string)
	for _, gateTypeItem := range gateTypes {
		gateTypesMap[gateTypeItem.ID] = gateTypeItem.Name
	}

	var wg sync.WaitGroup
	gateTypeInt, _ := strconv.ParseInt(gateType, 10, 32)
	switch gateTypesMap[int32(gateTypeInt)] {
	case "Промышленные ворота":
		wg.Go(func() {
			cfg.IndustrialDrives = []models.IndustrialGateDrive{}
			if err := database.DB.Find(&cfg.IndustrialDrives).Error; err != nil {
				helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
				return
			}
		})
	case "Бытовые ворота":
		wg.Go(func() {
			cfg.ResidentialDrives = []models.ResidentialGateDrive{}
			if err := database.DB.Find(&cfg.ResidentialDrives).Error; err != nil {
				helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
				return
			}
		})
	}

	wg.Go(func() {
		if err := database.DB.Find(&cfg.LiftTypes).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.Colors).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.Find(&cfg.CycleAmounts).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(width) as max_width, MIN(width) as min_width").
			Where("gate_type_id = ?", gateType).
			Scan(&cfg.WidthParams).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	})

	wg.Go(func() {
		if err := database.DB.
			Model(&models.Size{}).
			Select("MAX(height) as max_height, MIN(height) as min_height").
			Where("gate_type_id = ?", gateType).
			Scan(&cfg.HeightParams).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	})

	wg.Wait()
	data := map[string]any{
		"css":  "/assets/styles/calculator.css",
		"cfg":  cfg,
		"user": user,
	}

	if err := templates.ExecuteTemplate(w, "calculator.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
	}
}

// Route: /sizes
// Method: GET
func GetPriceBasedOnSize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width := query.Get("width")
	height := query.Get("height")
	gateType := query.Get("gateType")

	sessionToken, err := r.Cookie("session_token")
	var role string
	if err == nil {
		token := sessionToken.Value
		user := helpers.GetUserBySessionToken(token)
		role = user.Role.Name
	}

	var price int64
	if role == "dealer" || role == "manager" || role == "admin" {
		if err := database.DB.Model(&models.Size{}).
			Select("WholesalePrice").
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type_id = ?", gateType).
			Limit(1).
			Order("width asc, height asc").
			Pluck("WholesalePrice", &price).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		if err := database.DB.Model(&models.Size{}).
			Select("RetailPrice").
			Where("width >= ? AND height >= ?", width, height).
			Where("gate_type_id = ?", gateType).
			Limit(1).
			Order("width asc, height asc").
			Pluck("RetailPrice", &price).Error; err != nil {
			helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"price": price,
	})
}

// Route: /tables/{table_id}
// Method: GET
func GetDataBaseRedactor(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)
	role := user.Role.Name
	if role != "admin" {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	//migrator := database.DB.Migrator()

}

// Route: /tables
// Method: GET
func GetDataBaseTableList(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	token := sessionToken.Value
	user := helpers.GetUserBySessionToken(token)
	role := user.Role.Name
	if role != "admin" {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	migrator := database.DB.Migrator()
	tables, err := migrator.GetTables()
	if err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
	}

	data := map[string]any{
		"css":    "",
		"tables": tables,
		"user":   user,
	}

	if err := templates.ExecuteTemplate(w, "tables.html", data); err != nil {
		helpers.WriteErrorRelease(w, err, http.StatusInternalServerError)
	}
}
