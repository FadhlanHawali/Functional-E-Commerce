package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Order struct {
	Barang   Product_DB  `json:"barang"`
	Quantity int      `json:"quantity"`
	Total    int      `json:"total"`
	Status   string   `json:"status"`
	Customer CustomerRepo `json:"customer"`
}

type Cart struct {
	Id_Barang int `json:"idBarang"`
	Quantity int `json:"quantity"`
}
type OrderRepo struct {
	Id          int    `db:"id"`
	Id_Barang   int    `db:"id_barang" json:"idBarang"`
	Quantity    int    `db:"quantity" json:"quantity"`
	Status      string `db:"status" json:"status"`
	Total       int    `db:"total" json:"total"`
	TokenPayment	sql.NullString	`db:"token_payment"`
	Id_Customer int    `db:"id_customer" json:"idCustomer"`
	Id_Store	int    `db:"id_store"`
}

func (db *InDB) productAvailable(id_store int, id_barang int, qty int) bool {
	tx := db.DB.MustBegin()
	id := 0
	tx.Get(&id, fmt.Sprintf("SELECT id FROM products WHERE id = %d AND id_store = %d AND quantity >= %d", id_barang, id_store, qty))
	if err := tx.Commit(); err != nil {
		return false
	}
	if id > 0 {
		return true
	}
	return false
}

func (db *InDB) CreateAndListOrder(w http.ResponseWriter, r *http.Request) {
	var id_store int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["store_id"].(float64)
		id_store = int(tempId)
	} else {
		utils.WrapAPIError(w, r, "invalid token", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		CreateOrder(w, r, db, id_store)
		return
	} else if r.Method == "GET" {
		ListOrder(w, r, db, id_store)
		return
	}
	utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	return
}

func CreateOrder(w http.ResponseWriter, r *http.Request, db *InDB, id_store int) {
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var newOrder OrderRepo
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
		return
	}
	if !db.productAvailable(id_store, newOrder.Id_Barang, newOrder.Quantity) {
		utils.WrapAPIError(w, r, "Product is not available", http.StatusBadRequest)
		return
	}
	tx := db.DB.MustBegin()

	var barang Product_DB
	tx.Get(&barang, "SELECT * FROM products WHERE id = ? AND id_store = ?", newOrder.Id_Barang, id_store)
	total := newOrder.Quantity * barang.Price
	tx.MustExec("INSERT INTO orders (id_barang, id_customer, quantity, total, status, id_store) VALUES (?, ?, ?, ?, ?, ?)", newOrder.Id_Barang, newOrder.Id_Customer, newOrder.Quantity, total, "1", id_store)
	tx.Get(&newOrder.Id, "SELECT LAST_INSERT_ID() as id")
	tx.MustExec("UPDATE products SET quantity = quantity - ? WHERE id = ? and quantity > 0", newOrder.Quantity, newOrder.Id_Barang)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error creating new order", http.StatusInternalServerError)
		return
	}

	result, err := CreatePayment(w, r, newOrder, db)
	if err != nil {
		utils.WrapAPIData(w, r, newOrder, http.StatusAccepted, "success")
		return
	}

	utils.WrapAPIData(w, r, result, http.StatusOK, "success")
	return
}

func ListOrder(w http.ResponseWriter, r *http.Request, db *InDB, id_store int) {
	if r.Method != "GET" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	orderList := []OrderRepo{}
	tx := db.DB.MustBegin()
	err := tx.Select(&orderList, "SELECT * FROM orders WHERE id_store = ?", id_store); if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit(); if err != nil {
		utils.WrapAPIError(w, r, "error getting product", http.StatusInternalServerError)
		return
	}
	response := make([]*Order, len(orderList))
	for i, item := range orderList {
		var barang Product_DB
		err := db.DB.Get(&barang, "SELECT * FROM products WHERE id = ?", item.Id_Barang); if err != nil {
			fmt.Println(err)
		}
		var customer CustomerRepo
		db.DB.Get(&customer, "SELECT * FROM customers WHERE id = ?", item.Id_Customer)
		response[i] = &Order{
			Barang:   barang,
			Quantity: item.Quantity,
			Total:    item.Total,
			Status:   item.Status,
			Customer: customer,
		}
	}
	utils.WrapAPIData(w, r, response, http.StatusOK, "success")
	return
}

func (db *InDB) OrderController(w http.ResponseWriter, r *http.Request) {
	var id_store int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["store_id"].(float64)
		id_store = int(tempId)
	} else {
		utils.WrapAPIError(w, r, "invalid token", http.StatusBadRequest)
		return
	}

	id_order, err := strconv.Atoi(mux.Vars(r)["order"])
	if err != nil {
		utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		GetOrder(w, r, db, id_store, id_order)
		return
	} else if r.Method == "UPDATE" {
		UpdateOrder(w, r, db, id_store, id_order)
		return
	} else if r.Method == "DELETE" {
		DeleteOrder(w, r, db, id_store, id_order)
		return
	}
	utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	return
}

func GetOrder(w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_order int) {
	if r.Method != "GET" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var order OrderRepo
	tx := db.DB.MustBegin()
	tx.Select(&order, "SELECT * FROM orders WHERE id = ? AND id_store = ?", id_order, id_store)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error getting product", http.StatusInternalServerError)
		return
	}

	var barang Product_DB
	db.DB.Get(&barang, "SELECT * FROM products WHERE id = ? AND id_store = ?", order.Id_Barang, id_store)
	var customer CustomerRepo
	db.DB.Get(&customer, "SELECT * FROM customers WHERE id = ? AND id_store = ?", order.Id_Customer, id_store)

	response := &Order{
		Barang:   barang,
		Quantity: order.Quantity,
		Total:    order.Total,
		Status:   order.Status,
		Customer: customer,
	}
	utils.WrapAPIData(w, r, response, http.StatusOK, "success")
	return
}

func UpdateOrder(w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_order int) {

	if r.Method != "PUT" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var order OrderRepo
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.WrapAPIError(w, r, "Can't decode request body", http.StatusBadRequest)
		return
	}
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE orders SET status = ? WHERE id = ? AND id_store = ?", order.Status, id_order, id_store)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error updating order status", http.StatusInternalServerError)
		return
	}
	utils.WrapAPIData(w, r, order, http.StatusOK, "success")
	return
}

func DeleteOrder(w http.ResponseWriter, r *http.Request, db *InDB, id_store int, id_order int) {
	if r.Method != "DELETE" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var order OrderRepo
	tx := db.DB.MustBegin()
	tx.Select(&order, "SELECT * FROM orders WHERE id = ? AND id_store = ?", id_order, id_store)
	if order.Status == "1" {
		tx.MustExec("UPDATE products SET quantity = quantity + ? WHERE id = ? AND id_store = ?", order.Quantity, order.Id_Barang, id_store)
	}
	tx.MustExec("DELETE from orders where id = ? AND id_barang = ? AND id_store = ?", id_order, order.Id_Barang, id_store)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error delete order", http.StatusInternalServerError)
		return
	}
	utils.WrapAPISuccess(w, r, "success deleting order", http.StatusOK)
	return
}
