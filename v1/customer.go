package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
)

type Customer struct {
	Id string
	Email string `json:"email"`
	Nama string `json:"nama"`
	AlamatPengiriman string `json:"alamatPengiriman"`
}

type CustomerRepo struct {
	Id int `db:"id" json:"id"`
	Nama string `db:"cust_name" json:"nama"`
	Email string `db:"cust_email" json:"email"`
	AlamatPengiriman string `db:"cust_address" json:"alamat"`
	IdStore int `db:"id_store"`
}

func (db *InDB) rowExists(query string, args ...interface{}) bool {
    var exists bool
    query = fmt.Sprintf("SELECT exists (%s)", query)
    err := db.DB.QueryRow(query, args...).Scan(&exists)
    if err != nil {
        log.Fatalf("error checking if row exists '%s' %v", args, err)
    }
    return exists
}

func (db *InDB) CreateAndListCustomer (w http.ResponseWriter, r *http.Request) {
	var id_store int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["store_id"].(float64)
		id_store = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		AddCustomer(w, r, db, id_store)
	} else if r.Method == "GET" {
		ListCustomer(w, r, db, id_store)
	}
}

func AddCustomer (w http.ResponseWriter, r *http.Request, db *InDB, id_store int) {
	if r.Method != "POST" {
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	id := 0
	tx := db.DB.MustBegin()
	tx.Get(&id, "SELECT id FROM customers WHERE cust_email = ? AND id_store = ?", customer.Email, id_store)
	if id > 0 {
		tx.Select(&customer, "SELECT * FROM customers WHERE id = ?", id)
		utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
		return
	}

	tx.MustExec("INSERT INTO customers (cust_name, cust_address, cust_email, id_store) VALUES (?, ?, ?, ?)", customer.Nama, customer.AlamatPengiriman, customer.Email, id_store)
	tx.Get(&customer.Id, "SELECT LAST_INSERT_ID() as id")
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error inserting new customer", http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
	return
}

func ListCustomer (w http.ResponseWriter, r *http.Request, db *InDB, id_store int) {

	if r.Method != "GET"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}
	var customers []CustomerRepo

	tx := db.DB.MustBegin()
	tx.Select(&customers, "SELECT * FROM PRODUCTS WHERE id_store = ?", id_store)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error getting product",http.StatusInternalServerError)
		return
	}
	utils.WrapAPIData(w,r, customers, http.StatusOK, "success")
}

func (db *InDB) CustomerController (w http.ResponseWriter, r *http.Request) {
	var id_store int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["store_id"].(float64)
		id_store = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	id_customer, err := strconv.Atoi(mux.Vars(r)["customer"]); if err != nil {
		utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		GetCustomer(w, r, db, id_store, id_customer)
	} else if r.Method == "PUT" {
		UpdateCustomer(w, r, db, id_store, id_customer)
	} else if r.Method == "DELETE" {
		DeleteCustomer(w, r, db, id_store, id_customer)
	}
	utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
	return
}

func GetCustomer (w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_customer int) {
	if r.Method != "GET"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	var customer CustomerRepo
	tx := db.DB.MustBegin()
	tx.Select(&customer, "SELECT * FROM customers WHERE id =  ? AND id_store", id_customer, id_store)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error get user",http.StatusInternalServerError)
		return
	}
	utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
	return
}

func UpdateCustomer (w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_customer int) {
	if r.Method != "PUT"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE customers SET cust_name = ?, cust_address = ?, cust_email = ? WHERE id = ? AND id_store = ?", customer.Nama, customer.AlamatPengiriman, customer.Email, id_customer, id_store)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error updating user",http.StatusInternalServerError)
		return
	}
	utils.WrapAPIData(w, r, customer, http.StatusOK, "success")
	return
}

func DeleteCustomer (w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_customer int) {
	if r.Method != "DELETE"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	tx := db.DB.MustBegin()
	tx.MustExec("DELETE from customers where id = ? AND id_store = ?" , id_customer, id_store)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error delete product",http.StatusInternalServerError)
		return
	}
	utils.WrapAPISuccess(w,r,"success deleting user",http.StatusOK)
	return

}