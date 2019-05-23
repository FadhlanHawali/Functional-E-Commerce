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

func (db *InDB) AddProduct (w http.ResponseWriter, r *http.Request){

	var product Product
	var id int

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product);err!=nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id = int(tempId)

		//email = tokenMap["email"].(string)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM STORES WHERE id_user = %d",id))

	tx.MustExec("INSERT INTO products (prod_name,quantity,description,price,url_pic,id_store) VALUES (?, ? ,? , ?, ?, ?)",product.NamaBarang,product.Quantity,product.Deskripsi,product.Harga,product.UrlGambar,id)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error adding product",http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w,r,"success adding product",http.StatusCreated)

}

func (db *InDB) ListProduct (w http.ResponseWriter, r *http.Request){
	var product []Product_DB
	var id int
	if r.Method != "GET"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id = int(tempId)

		//email = tokenMap["email"].(string)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM STORES WHERE id_user = %d",id))
	tx.Select(&product,fmt.Sprintf("SELECT * FROM PRODUCTS WHERE id_store = %d",id))
	//tx.MustExec("INSERT INTO products (prod_name,quantity,description,price,url_pic,id_store) VALUES (?, ? ,? , ?, ?, ?)",product.NamaBarang,product.Quantity,product.Deskripsi,product.Harga,product.UrlGambar,id)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error getting product",http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w,r,product,http.StatusOK,"success")
}

func (db *InDB) DeleteProduct (w http.ResponseWriter, r *http.Request){
	var id int
	var id_product int
	var err error
	if r.Method != "DELETE"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	id_product_temp := mux.Vars(r)["id"]
	if id_product,err = strconv.Atoi(id_product_temp);err != nil{
		utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
		return
	}
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.Get(&id, fmt.Sprintf("SELECT id FROM STORES WHERE id_user = %d",id))
	tx.MustExec(fmt.Sprintf("DELETE from products where id_store = %d AND id = %d",id,id_product))
	//tx.MustExec("INSERT INTO products (prod_name,quantity,description,price,url_pic,id_store) VALUES (?, ? ,? , ?, ?, ?)",product.NamaBarang,product.Quantity,product.Deskripsi,product.Harga,product.UrlGambar,id)

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
