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
type Order struct {
	Barang Product `json:"barang"`
	Quantity int `json:"quantity"`
	Total int `json:"total"`
	Status string `json:"status"`
	Customer Customer `json:"customer"`
}

type OrderRepo struct {
	Id int `db:"id"`
	Id_Barang int `db:"id_barang" json:"idBarang"`
	Quantity int `db:"quantity" json:"quantity"`
	Status string `db:"status" json:"status"`
	Total int `db:"total" json:"total"`
	Id_Customer int `db:"id_customer" json:"idCustomer"`
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

func (db *InDB) CreateOrder (w http.ResponseWriter, r *http.Request) {
	id_store, err := strconv.Atoi(mux.Vars(r)["id"]); if err != nil {
		utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
		return
	}
	var id_user int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id_user = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}
	if (!db.isMyStore(id_user, id_store)) {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	if r.Method != "POST" {
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}
	var newOrder OrderRepo
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}
	if (!db.productAvailable(id_store, newOrder.Id_Barang, newOrder.Quantity)) {
		utils.WrapAPIError(w,r,"Product is not available",http.StatusBadRequest)
		return
	}
	tx := db.DB.MustBegin()

	var barang Product_DB
	tx.Get(&barang, "SELECT * FROM products WHERE id = ? AND id_store = ?", newOrder.Id_Barang, id_store)
	total := newOrder.Quantity * barang.Price
	tx.MustExec("INSERT INTO orders (id_barang, id_customer, quantity, total, status) VALUES (?, ?, ?, ?, ?)", newOrder.Id_Barang, newOrder.Id_Customer, newOrder.Quantity, total, "1")
	tx.Get(&newOrder.Id, "SELECT LAST_INSERT_ID() as id")
	tx.MustExec("UPDATE products SET quantity = quantity - ? WHERE id = ? and quantity > 0", newOrder.Quantity, newOrder.Id_Barang)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error creating new order", http.StatusInternalServerError)
		return
	}

	result,err := CreatePayment(w,r,newOrder,db);if err!=nil{
		utils.WrapAPIData(w, r, newOrder, http.StatusAccepted, "success")
		return
	}

	utils.WrapAPIData(w,r,result,http.StatusOK,"success")
	return
}

func (db *InDB) OrderController (w http.ResponseWriter, r *http.Request) {
	id_store, err := strconv.Atoi(mux.Vars(r)["id"]); if err != nil {
		utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
		return
	}
	var id_user int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id_user = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}
	if (!db.isMyStore(id_user, id_store)) {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		var orders []OrderRepo
		tx := db.DB.MustBegin()
		tx.Select(&orders,fmt.Sprintf("SELECT * FROM orders WHERE id_store = %d", id_store))
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w,r,"error getting product",http.StatusInternalServerError)
			return
		}
		response := make([]*Order, len(orders))
		for i, item := range orders {
			var barang Product
			db.DB.Select(&barang, "SELECT * FROM products WHERE id = ?", item.Id_Barang)
			var customer Customer
			db.DB.Select(&customer, "SELECT * FROM customers WHERE id = ?", item.Id_Customer)
			response[i] = &Order{
				Barang: barang,
				Quantity: item.Quantity,
				Total: item.Total,
				Status: item.Status,
				Customer: customer,
			}
		}
		utils.WrapAPIData(w, r, response, http.StatusOK, "success")
		return
	} else if r.Method == "UPDATE" {
		id_order, err := strconv.Atoi(mux.Vars(r)["order"]); if err != nil {
			utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
			return
		}
		var order OrderRepo
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
			return
		}
		tx := db.DB.MustBegin()
		tx.MustExec("UPDATE orders SET status = ? WHERE id = ?", order.Status, id_order)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w,r,"error updating order status",http.StatusInternalServerError)
			return
		}
		utils.WrapAPIData(w, r, order, http.StatusOK, "success")
		return

	} else if r.Method == "DELETE" {
		id_order, err := strconv.Atoi(mux.Vars(r)["order"]); if err != nil{
			utils.WrapAPIError(w,r,"error converting string to integer",http.StatusInternalServerError)
			return
		}
		var order OrderRepo
		tx := db.DB.MustBegin()
		tx.Select(&order, "SELECT * FROM orders WHERE id = ?", id_order)
		if (order.Status == "1") {
			tx.MustExec("UPDATE products SET quantity = quantity + ? WHERE id = ?", order.Quantity, order.Id_Barang)
		}
		tx.MustExec("DELETE from orders where id = ? AND id_barang = ?", id_order, order.Id_Barang)
		if err := tx.Commit(); err != nil {
			utils.WrapAPIError(w,r,"error delete order",http.StatusInternalServerError)
			return
		}
		utils.WrapAPISuccess(w,r,"success deleting order",http.StatusOK)
		return
	}
	utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
	return
}
