package main

import (
	"github.com/gorilla/mux"
	"github.com/FadhlanHawali/Functional-E-Commerce/v1"
	"log"
	"fmt"
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/database"
)

func main(){
	conn, err := database.InitDb("root:pintar123@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
		return
	}
	defer conn.DB.Close()
	api := &v1.InDB{DB: conn.GetDB()}
	router := mux.NewRouter()
	//TODO API APA AJA
	router.HandleFunc("/api/v1/create", api.CreateLapak)
	//TODO
	http.Handle("/", router)
	port := ":8080"
	log.Printf("Server Running on port %s",port)
	log.Fatal(http.ListenAndServe(port, router))
}
