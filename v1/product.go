package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
)

func (db *InDB) CreateAndListProduct (w http.ResponseWriter, r *http.Request) {
	var id_store int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id_store"].(float64)
		id_store = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
		AddProduct(w, r, db, id_store)
	} else if (r.Method == "GET") {
		ListProduct(w, r, db, id_store)
	}
}

func (db *InDB) ProductController (w http.ResponseWriter, r *http.Request) {
	var id_store int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id_store"].(float64)
		id_store = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	id_product, err := strconv.Atoi(mux.Vars(r)["product"]); if err != nil {
		utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
		return
	}

	if r.Method == "DELETE" {
		DeleteProduct(w, r, db, id_store, id_product)
	}
}

func AddProduct (w http.ResponseWriter, r *http.Request, db *InDB, id_store int) {

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}
	var product Product

	if err := json.NewDecoder(r.Body).Decode(&product);err!=nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO products (prod_name,quantity,description,price,url_pic,id_store) VALUES (?, ? ,? , ?, ?, ?)", product.NamaBarang, product.Quantity, product.Deskripsi, product.Harga, product.UrlGambar, id_store)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error adding product",http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w,r,"success adding product",http.StatusCreated)
}

func ListProduct (w http.ResponseWriter, r *http.Request, db *InDB, id_store int)  {
	if r.Method != "GET"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	var product []Product_DB

	tx := db.DB.MustBegin()
	tx.Select(&product, "SELECT * FROM products WHERE id_store = ?",id_store)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error getting product",http.StatusInternalServerError)
		return
	}
	utils.WrapAPIData(w,r,product,http.StatusOK,"success")
}

func DeleteProduct (w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_product int) {
	if r.Method != "DELETE"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	tx := db.DB.MustBegin()
	tx.MustExec(fmt.Sprintf("DELETE from products where id_store = %d AND id = %d", id_store, id_product))

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error getting product",http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w,r,"success deleting product",http.StatusOK)
}

type Product struct {
	NamaBarang string `json:"namaBarang"`
	Deskripsi string `json:"deskripsi"`
	Quantity int `json:"quantity"`
	Harga int `json:"harga"`
	UrlGambar string `json:"urlGambar"`
}

type Product_DB struct {
	Id int `db:"id"`
	Prod_Name string `db:"prod_name"`
	Quantity int `db:"quantity"`
	Description string `db:"description"`
	Price int `db:"price"`
	Url_Pic string `db:"url_pic"`
	Id_Store int `db:"id_store"`
}
