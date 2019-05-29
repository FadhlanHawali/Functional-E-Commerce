package v1

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const TokenContextKey = "MyAppToken"


func (db *InDB) CreateAndListLapak(w http.ResponseWriter, r *http.Request){
	var id_user int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id_user = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}
	if r.Method == "POST"{
		CreateLapak(w, r, db, id_user)
		return
	} else if r.Method == "GET" {
		ListLapak(w, r, db, id_user)
		return
	}
	utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
	return
}

func ListLapak (w http.ResponseWriter, r *http.Request, db *InDB, id_user int) {
	if r.Method != "GET" {
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}
	lapaks := []Lapak{}
	err := db.DB.Select(&lapaks, "SELECT id, store_name, address, handphone, bank_number FROM stores WHERE id_user = ?", id_user); if err != nil {
		fmt.Println(err)
		utils.WrapAPIError(w, r, "error getting lapak", http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, lapaks, http.StatusOK, "success")
	return

}

func CreateLapak (w http.ResponseWriter, r *http.Request, db *InDB, id_user int) {
	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	var lapak Lapak
	if err := json.NewDecoder(r.Body).Decode(&lapak);err!=nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO stores (store_name,address,handphone,bank_number,id_user) VALUES (?, ? ,? , ?, ?)",lapak.NamaLapak,lapak.Alamat,lapak.Telepon,lapak.NomerRekening,id_user)
	var id int
	tx.Get(&id, "SELECT LAST_INSERT_ID() as id")
	jwt := jwt.MapClaims{
		"store_id": id,
		"store_name": lapak.NamaLapak,
	}
	apiKey, err := utils.GenerateToken(w, r, jwt, "storesecr3t"); if err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}
	lapak.ApiKey = apiKey
	tx.MustExec("UPDATE stores SET api_key=? WHERE id=?", lapak.ApiKey, id)
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error creating new lapak",http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, map[string]interface{} {
		"Nama": lapak.NamaLapak,
		"ApiKey": lapak.ApiKey,
	}, http.StatusOK, "success")
}

func (db *InDB) StoreController (w http.ResponseWriter, r *http.Request) {
	var id_user int
	if token := r.Context().Value(TokenContextKey); token != nil {
		tokenMap := token.(jwt.MapClaims)
		tempId := tokenMap["id"].(float64)
		id_user = int(tempId)
	} else {
		utils.WrapAPIError(w,r,"invalid token",http.StatusBadRequest)
		return
	}

	id_store, err := strconv.Atoi(mux.Vars(r)["store"])
	if err != nil {
		utils.WrapAPIError(w, r, "error converting string to integer", http.StatusInternalServerError)
		return
	}

	if r.Method == "GET"{
		UpdateApiKey(w, r, db, id_user, id_store)
		return
	} else if r.Method == "DELETE" {
		DeleteLapak(w, r, db, id_user, id_store)
		return
	}
}

func UpdateApiKey(w http.ResponseWriter, r *http.Request, db *InDB, id_user int, id_store int) {
	if r.Method != "GET" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var lapak Lapak
	tx := db.DB.MustBegin()
	tx.Get(&lapak, "SELECT * FROM stores WHERE id_user = ? AND id = ?", id_user, id_store)
	jwt := jwt.MapClaims{
		"store_id": id_store,
		"store_name": lapak.NamaLapak,
	}
	apiKey, err := utils.GenerateToken(w, r, jwt, "storesecr3t"); if err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}
	lapak.ApiKey = apiKey
	tx.MustExec("UPDATE stores SET api_key=? WHERE id=?", lapak.ApiKey, id_store)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		utils.WrapAPIError(w,r,"error creating new lapak",http.StatusInternalServerError)
		return
	}
	utils.WrapAPIData(w, r, map[string]interface{} {
		"Nama": lapak.NamaLapak,
		"ApiKey": lapak.ApiKey,
	}, http.StatusOK, "success")
}

func DeleteLapak (w http.ResponseWriter, r *http.Request, db *InDB, id_user int, id_store int) {
	if r.Method != "DELETE" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	tx := db.DB.MustBegin()
	tx.MustExec("DELETE FROM stores where id = ? AND id_user = ?", id_store, id_user)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		utils.WrapAPIError(w, r, "error delete order", http.StatusInternalServerError)
		return
	}
	utils.WrapAPISuccess(w, r, "success deleting order", http.StatusOK)
	return
}


type Lapak struct {
	Id int `json: "id"`
	NamaLapak string `db:"store_name" json:"namaLapak"`
	Alamat string `db:"address" json:"alamat"`
	Telepon string `db:"handphone" json:"telepon"`
	NomerRekening string `db:"bank_number" json:"nomerRekening"`
	ApiKey string `db:"api_key" json:"apiKey"`
}
