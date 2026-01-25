package main

import (
	"log"
	"net/http"
	"project/pkg/database"
	"project/pkg/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.InitDB()
	database.AutoMigrateAll(db)

	r := mux.NewRouter()

	http.Handle("/assets/scripts/",
		http.StripPrefix("/assets/scripts/",
			http.FileServer(http.Dir("./assets/scripts"))))

	http.Handle("/assets/images/",
		http.StripPrefix("/assets/images/",
			http.FileServer(http.Dir("./assets/images"))))

	http.Handle("/assets/styles/",
		http.StripPrefix("/assets/styles/",
			http.FileServer(http.Dir("./assets/styles"))))

	// Главная страница
	r.HandleFunc("/", handlers.MainHandler).Methods("GET")

	// Контакты
	r.HandleFunc("/contacts", handlers.ContactsHandler).Methods("GET")

	// Аутентификация
	r.HandleFunc("/log", handlers.GetSignInForm).Methods("GET")
	r.HandleFunc("/log", handlers.SignIn).Methods("POST")
	r.HandleFunc("/reg", handlers.GetSignUpFrom).Methods("GET")
	r.HandleFunc("/reg", handlers.SignUp).Methods("POST")
	r.HandleFunc("/log_out", handlers.LogOut).Methods("POST")

	// Аккаунт
	r.HandleFunc("/user", handlers.GetUserInfo).Methods("GET")
	r.HandleFunc("/user/dealers", handlers.GetUserDealers).Methods("GET")

	// Заказы
	r.HandleFunc("/orders", handlers.GetAllUserOrders).Methods("GET")
	r.HandleFunc("/orders", handlers.CreateNewOrder).Methods("POST")

	r.HandleFunc("/orders/{order_id}", handlers.GetUserOrderById).Methods("GET")
	r.HandleFunc("/orders/{order_id}", handlers.AddNewGateInOrder).Methods("POST")
	r.HandleFunc("/orders/{order_id}", handlers.DeleteUserOrder).Methods("DELETE")

	r.HandleFunc("/orders/{order_id}/products", handlers.GetProductsInOrder).Methods("GET")
	r.HandleFunc("/orders/{order_id}/products", handlers.AddNewProductInOrder).Methods("POST")
	r.HandleFunc("/orders/{order_id}/products/{product_id}", handlers.UpdateProductList).Methods("PUT")
	r.HandleFunc("/orders/{order_id}/products/{product_id}", handlers.DeleteProductFromOrder).Methods("DELETE")

	r.HandleFunc("/orders/{order_id}/{gate_id}", handlers.GetGateInOrder).Methods("GET")
	r.HandleFunc("/orders/{order_id}/{gate_id}", handlers.DeleteGateFromOrder).Methods("DELETE")
	r.HandleFunc("/orders/{order_id}/{gate_id}", handlers.UpdateGateInOrder).Methods("PUT")

	// Доделать
	r.HandleFunc("/orders/{order_id}/{gate_id}/options", handlers.GetGateOptions).Methods("GET")
	r.HandleFunc("/orders/{order_id}/documents", handlers.GetOrderDocuments).Methods("GET")

	r.HandleFunc("/calculator", handlers.GetCalculatorForUser).Methods("GET")
	r.HandleFunc("/gate_type", handlers.GetGateTypesList).Methods("GET")

	r.HandleFunc("/sizes", handlers.GetPriceBasedOnSize).Methods("GET")

	// Администратор
	r.HandleFunc("/tables/{table_name}", handlers.GetDataBaseRedactor).Methods("GET")
	r.HandleFunc("/tables", handlers.GetDataBaseTableList).Methods("GET")

	http.Handle("/", r)
	err_start := http.ListenAndServe(":8080", nil)
	log.Fatal(err_start)
}
