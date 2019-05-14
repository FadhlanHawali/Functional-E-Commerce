package main

import (
	"github.com/gorilla/mux"
	"github.com/FadhlanHawali/Functional-E-Commerce/v1"
	"log"
	"net/http"
)

func main(){

	router := mux.NewRouter()
	//TODO API APA AJA
	router.HandleFunc("/api/v1/create", v1.CreateLapak)
	//TODO
	http.Handle("/", router)
	port := ":8080"
	log.Printf("Server Running on port %s",port)
	log.Fatal(http.ListenAndServe(port, router))
}
