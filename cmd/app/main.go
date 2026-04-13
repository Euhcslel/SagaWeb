package main

import (
	"log"
	"net/http"
	"os"
	"project/internal/database"
	"project/internal/helpers"
	handlers "project/internal/transport/http"

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

	mux := http.NewServeMux()

	assets := http.StripPrefix(
		"/web/assets/",
		http.FileServer(http.Dir("./web/assets")),
	)

	mux.Handle("GET /web/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		assets.ServeHTTP(w, r)
	}))

	protoAssets := http.StripPrefix(
		"/api/proto/",
		http.FileServer(http.Dir("./api/proto")),
	)

	mux.Handle("GET /api/proto/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		protoAssets.ServeHTTP(w, r)
	}))

	// Главная страница
	mux.HandleFunc("GET /", handlers.MainHandler)

	// Аутентификация
	mux.HandleFunc("GET /sign_in", handlers.SignInForm)
	mux.HandleFunc("POST /sign_in", handlers.SignIn)
	mux.HandleFunc("GET /sign_up", handlers.SignUpForm)
	mux.HandleFunc("POST /sign_up", handlers.SignUp)
	mux.HandleFunc("POST /sign_out", handlers.SignOut)

	// Аккаунт
	mux.HandleFunc("GET /user", handlers.GetUserInfo)
	mux.HandleFunc("GET /user/dealers", handlers.GetUserDealers)

	// Заказы
	mux.HandleFunc("GET /orders", handlers.GetAllUserOrders)
	mux.HandleFunc("POST /orders", handlers.CreateNewOrder)

	mux.HandleFunc("GET /orders/{order_id}", handlers.GetUserOrderById)
	mux.HandleFunc("POST /orders/{order_id}", handlers.AddNewGateInOrder)
	mux.HandleFunc("DELETE /orders/{order_id}", handlers.DeleteUserOrder)

	mux.HandleFunc("POST /orders/{order_id}/products", handlers.AddNewProductInOrder)
	mux.HandleFunc("PUT /orders/{order_id}/products/{product_id}", handlers.UpdateProductList)
	mux.HandleFunc("DELETE /orders/{order_id}/products/{product_id}", handlers.DeleteProductFromOrder)

	mux.HandleFunc("GET /orders/{order_id}/{gate_id}", handlers.GetGateInOrder)
	mux.HandleFunc("DELETE /orders/{order_id}/{gate_id}", handlers.DeleteGateFromOrder)
	mux.HandleFunc("PUT /orders/{order_id}/{gate_id}", handlers.UpdateGateInOrder)

	// Доделать
	mux.HandleFunc("GET /orders/{order_id}/documents", handlers.GetOrderDocuments)

	mux.HandleFunc("GET /calculator", handlers.GetCalculatorForUser)

	mux.HandleFunc("GET /sizes", handlers.GetPriceBasedOnSize)

	// Администратор
	mux.HandleFunc("GET /tables/{table_name}", handlers.GetDataBaseRedactor)
	mux.HandleFunc("GET /tables", handlers.GetDataBaseTableList)

	err_start := http.ListenAndServe(":8080", http.NewCrossOriginProtection().Handler(mux))
	log.Fatal(err_start)
}
