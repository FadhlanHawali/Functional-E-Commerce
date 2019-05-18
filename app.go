package main

import (
	"github.com/gorilla/mux"
	"github.com/FadhlanHawali/Functional-E-Commerce/v1"
	"log"
	"fmt"
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/database"
	"context"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
)

const TokenContextKey = "MyAppToken"

func main(){
	conn, err := database.InitDb("root:pintar123@tcp(127.0.0.1:3306)/commerce")
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
		return
	}
	defer conn.DB.Close()
	api := &v1.InDB{DB: conn.GetDB()}
	router := mux.NewRouter()
	//TODO API APA AJA
	router.HandleFunc("/api/v1/create", WithAuth(http.HandlerFunc(api.CreateLapak)))
	router.HandleFunc("/api/v1/user/create",api.CreateUser)
	//TODO
	http.Handle("/", router)
	port := ":8080"
	log.Printf("Server Running on port %s",port)
	log.Fatal(http.ListenAndServe(port, router))
}

func WithAuth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			utils.WrapAPIError(w,r,"Authorization Header can't be Empty",http.StatusUnauthorized) // continue without token
			return
		}

		token, err := utils.ValidateToken(auth)
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