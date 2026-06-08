package main

import (
	"github.com/Euhcslel/SagaWeb/internal/database"
	handlers "github.com/Euhcslel/SagaWeb/internal/transport/http"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err_load := godotenv.Load()
	if err_load != nil {
		log.Fatal("error loading .env file")
	}

	logFile, err := os.OpenFile("logs/error.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	database.InitDB()

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

	productImages := http.StripPrefix(
		"/data/images/",
		http.FileServer(http.Dir("./data/images")),
	)

	mux.Handle("GET /data/images/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		productImages.ServeHTTP(w, r)
	}))

	// Главная страница
	mux.Handle("GET /", handlers.WithOptionalAuth(http.HandlerFunc(handlers.MainHandler)))

	// Аутентификация
	mux.Handle("GET /sign_in", handlers.CacheControl("public, max-age=3600")(handlers.WithOptionalAuth(http.HandlerFunc(handlers.SignInForm))))
	mux.Handle("POST /sign_in", handlers.WithOptionalAuth(handlers.RateLimit(http.HandlerFunc(handlers.SignIn))))
	mux.Handle("GET /sign_up", handlers.CacheControl("public, max-age=3600")(handlers.WithOptionalAuth(http.HandlerFunc(handlers.SignUpForm))))
	mux.Handle("POST /sign_up", handlers.WithOptionalAuth(handlers.RateLimit(http.HandlerFunc(handlers.SignUp))))
	mux.Handle("POST /sign_out", handlers.WithOptionalAuth(http.HandlerFunc(handlers.SignOut)))

	// Аккаунт
	mux.Handle("GET /user", handlers.RequireAuth(http.HandlerFunc(handlers.GetUserInfo)))
	mux.Handle("POST /user", handlers.RequireAuth(http.HandlerFunc(handlers.UpdateUserInfo)))

	mux.Handle("GET /user/dealers", handlers.RequireAuth(http.HandlerFunc(handlers.GetUserDealers)))
	mux.Handle("POST /user/dealers/{dealer_id}/contract", handlers.RequireAuth(http.HandlerFunc(handlers.AttachContractToDealer)))

	// Заявки на регистрацию
	mux.Handle("GET /dealers/requests", handlers.RequireAuth(http.HandlerFunc(handlers.GetDealersRegRequests)))
	mux.Handle("POST /dealers/requests/{request_id}/confirm", handlers.RequireAuth(http.HandlerFunc(handlers.ConfirmDealerRegRequest)))
	mux.Handle("POST /dealers/requests/{request_id}/reject", handlers.RequireAuth(http.HandlerFunc(handlers.RejectDealerRegRequest)))

	// Заказы
	mux.Handle("GET /orders", handlers.RequireAuth(http.HandlerFunc(handlers.GetAllUserOrders)))
	mux.Handle("POST /orders", handlers.RequireAuth(http.HandlerFunc(handlers.CreateNewOrder)))

	mux.Handle("GET /orders/{order_id}", handlers.RequireAuth(http.HandlerFunc(handlers.GetUserOrderById)))
	mux.Handle("POST /orders/{order_id}", handlers.RequireAuth(http.HandlerFunc(handlers.AddNewGateInOrder)))
	mux.Handle("DELETE /orders/{order_id}", handlers.RequireAuth(http.HandlerFunc(handlers.DeleteUserOrder)))

	mux.Handle("PUT /orders/{order_id}/products", handlers.RequireAuth(http.HandlerFunc(handlers.UpdateProductList)))

	mux.Handle("PUT /orders/{order_id}/status", handlers.RequireAuth(http.HandlerFunc(handlers.UpdateOrderStatus)))

	mux.Handle("GET /orders/{order_id}/{gate_id}", handlers.RequireAuth(http.HandlerFunc(handlers.GetGateInOrder)))
	mux.Handle("DELETE /orders/{order_id}/{gate_id}", handlers.RequireAuth(http.HandlerFunc(handlers.DeleteGateFromOrder)))
	mux.Handle("PUT /orders/{order_id}/{gate_id}", handlers.RequireAuth(http.HandlerFunc(handlers.UpdateGateInOrder)))

	// Работа с документами
	mux.Handle("GET /orders/{order_id}/documents", handlers.RequireAuth(http.HandlerFunc(handlers.GetOrderDocuments)))
	mux.Handle("POST /orders/{order_id}/offer", handlers.RequireAuth(http.HandlerFunc(handlers.UploadOfferToOrder)))
	mux.Handle("POST /orders/{order_id}/bill", handlers.RequireAuth(http.HandlerFunc(handlers.UploadBillToOrder)))
	mux.Handle("POST /orders/{order_id}/contract", handlers.RequireAuth(http.HandlerFunc(handlers.UploadContractToOrder)))
	mux.Handle("GET /orders/{order_id}/documents/{document_type}/{document_name}", handlers.RequireAuth(http.HandlerFunc(handlers.DownloadOrderDocument)))
	mux.Handle("DELETE /orders/{order_id}/documents/{document_type}/{document_name}", handlers.RequireAuth(http.HandlerFunc(handlers.DeleteOrderDocument)))
	mux.Handle("POST /offer", handlers.WithOptionalAuth(http.HandlerFunc(handlers.GetOfferForOrder)))
	mux.Handle("GET /orders/{order_id}/appendix", handlers.RequireAuth(http.HandlerFunc(handlers.GetAppendixForOrder)))

	mux.Handle("GET /calculator", handlers.CacheControl("private, max-age=300")(handlers.WithOptionalAuth(http.HandlerFunc(handlers.GetCalculatorForUser))))
	mux.Handle("GET /catalog", handlers.WithOptionalAuth(http.HandlerFunc(handlers.GetCatalog)))

	mux.Handle("GET /sizes", handlers.WithOptionalAuth(http.HandlerFunc(handlers.GetPriceBasedOnSize)))

	// Логистик
	mux.Handle("GET /logistician/size-import", handlers.RequireAuth(http.HandlerFunc(handlers.GetSizeImportPage)))
	mux.Handle("POST /logistician/size-import", handlers.RequireAuth(http.HandlerFunc(handlers.ImportSizes)))
	mux.Handle("PUT /logistician/orders/{order_id}", handlers.RequireAuth(http.HandlerFunc(handlers.UpdateLogisticianOrder)))

	// Администратор
	mux.Handle("GET /tables/{table_name}", handlers.RequireAuth(http.HandlerFunc(handlers.GetDataBaseRedactor)))
	mux.Handle("POST /tables/{table_name}", handlers.RequireAuth(http.HandlerFunc(handlers.AddNewDataBaseTableRow)))
	mux.Handle("PUT /tables/{table_name}", handlers.RequireAuth(http.HandlerFunc(handlers.UpdateRow)))
	mux.Handle("DELETE /tables/{table_name}/{row_id}", handlers.RequireAuth(http.HandlerFunc(handlers.DeleteRow)))

	mux.Handle("GET /tables", handlers.RequireAuth(http.HandlerFunc(handlers.GetDataBaseTableList)))

	err = http.ListenAndServe(":80", http.NewCrossOriginProtection().Handler(handlers.SecurityHeaders(mux)))
	log.Fatal(err)
}
