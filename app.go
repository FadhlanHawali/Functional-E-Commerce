package main

import (
	"log"
	"fmt"
	"context"
	"net/http"

	"github.com/FadhlanHawali/Functional-E-Commerce/database"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"github.com/FadhlanHawali/Functional-E-Commerce/v1"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/rs/cors"
)

const TokenContextKey = "MyAppToken"

func main() {
	viper.SetConfigFile("./config/dev.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	conn, err := database.InitDb(fmt.Sprintf("%s:%s@tcp(%s:%s)/commerce", viper.Get("db.username"), viper.Get("db.password"), viper.Get("db.host"), viper.Get("db.port")))
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
		return
	}
	defer conn.DB.Close()
	api := &v1.InDB{DB: conn.GetDB()}
	router := mux.NewRouter()
	//TODO API APA AJA
	router.HandleFunc("/api/v1/order", WithStoreAuth(http.HandlerFunc(api.CreateAndListOrder)))
	router.HandleFunc("/api/v1/order/{order}", WithStoreAuth(http.HandlerFunc(api.OrderController)))

	router.HandleFunc("/api/v1/customer", WithStoreAuth(http.HandlerFunc(api.CreateAndListCustomer)))
	router.HandleFunc("/api/v1/customer/{customer}", WithStoreAuth(http.HandlerFunc(api.CustomerController)))

	router.HandleFunc("/api/v1/product", WithStoreAuth(http.HandlerFunc(api.CreateAndListProduct)))
	router.HandleFunc("/api/v1/product/{product}", WithStoreAuth(http.HandlerFunc(api.ProductController)))

	router.HandleFunc("/api/v1/store", WithAuth(http.HandlerFunc(api.CreateLapak)))

	router.HandleFunc("/api/v1/user/login", api.Login)
	router.HandleFunc("/api/v1/user/create", api.CreateUser)
	router.HandleFunc("/api/v1/store/order/{idOrder}/user/{idCustomer}/payment/{token}", api.UpdatePayment)
	//TODO

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8000"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	port := fmt.Sprintf(":%s", viper.Get("host.port"))
	log.Printf("Server Running on port %s",port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func WithAuth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			utils.WrapAPIError(w, r, "Authorization Header can't be Empty", http.StatusUnauthorized) // continue without token
			return
		}

		token, err := utils.ValidateToken(auth, "secret")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithStoreAuth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			utils.WrapAPIError(w, r, "Authorization Header can't be Empty", http.StatusUnauthorized) // continue without token
			return
		}

		token, err := utils.ValidateToken(auth, "storesecr3t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//Untuk manggil token nya

//func Handle(w http.ResponseWriter, r *http.Request) {
//	if token := r.Context().Value(TokenContextKey); token != nil {
//		// User is logged in
//	} else {
//		// User is not logged in
//	}
//}
