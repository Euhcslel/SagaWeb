package main

import (
	"log"
	"net/http"
	"os"
	"project/pkg/database"
	"project/pkg/handlers"
	"project/pkg/helpers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	logFile, err := os.Create("logs/error.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	helpers.LogFile = logFile

	r := mux.NewRouter()

	assets := http.StripPrefix(
		"/assets/",
		http.FileServer(http.Dir("./assets")),
	)

	r.PathPrefix("/assets/").Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
			assets.ServeHTTP(w, r)
		}),
	)

	// Главная страница
	r.HandleFunc("/", handlers.MainHandler).Methods("GET")

	// Аутентификация
	r.HandleFunc("/sign_in", handlers.SignInForm).Methods("GET")
	r.HandleFunc("/sign_in", handlers.SignIn).Methods("POST")
	r.HandleFunc("/sign_up", handlers.SignUpForm).Methods("GET")
	r.HandleFunc("/sign_up", handlers.SignUp).Methods("POST")
	r.HandleFunc("/sign_out", handlers.SignOut).Methods("POST")

	// Аккаунт
	r.HandleFunc("/user", handlers.GetUserInfo).Methods("GET")
	r.HandleFunc("/user/dealers", handlers.GetUserDealers).Methods("GET")

	// Заказы
	r.HandleFunc("/orders", handlers.GetAllUserOrders).Methods("GET")
	r.HandleFunc("/orders", handlers.CreateNewOrder).Methods("POST")

	r.HandleFunc("/orders/{order_id}", handlers.GetUserOrderById).Methods("GET")
	r.HandleFunc("/orders/{order_id}", handlers.AddNewGateInOrder).Methods("POST")
	r.HandleFunc("/orders/{order_id}", handlers.DeleteUserOrder).Methods("DELETE")

	r.HandleFunc("/orders/{order_id}/products", handlers.AddNewProductInOrder).Methods("POST")
	r.HandleFunc("/orders/{order_id}/products/{product_id}", handlers.UpdateProductList).Methods("PUT")
	r.HandleFunc("/orders/{order_id}/products/{product_id}", handlers.DeleteProductFromOrder).Methods("DELETE")

	r.HandleFunc("/orders/{order_id}/{gate_id}", handlers.GetGateInOrder).Methods("GET")
	r.HandleFunc("/orders/{order_id}/{gate_id}", handlers.DeleteGateFromOrder).Methods("DELETE")
	r.HandleFunc("/orders/{order_id}/{gate_id}", handlers.UpdateGateInOrder).Methods("PUT")

	// Доделать
	r.HandleFunc("/orders/{order_id}/documents", handlers.GetOrderDocuments).Methods("GET")

	r.HandleFunc("/calculator", handlers.GetCalculatorForUser).Methods("GET")
	r.HandleFunc("/gate_type", handlers.GetGateTypesList).Methods("GET")

	r.HandleFunc("/sizes", handlers.GetPriceBasedOnSize).Methods("GET")

	// Администратор
	r.HandleFunc("/tables/{table_name}", handlers.GetDataBaseRedactor).Methods("GET")
	r.HandleFunc("/tables", handlers.GetDataBaseTableList).Methods("GET")

	err_start := http.ListenAndServe(":8080", http.NewCrossOriginProtection().Handler(r))
	log.Fatal(err_start)
}
