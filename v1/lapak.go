package v1

import (
	"fmt"
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
)

const TokenContextKey = "MyAppToken"

func (db *InDB) CreateLapak(w http.ResponseWriter, r *http.Request){

	var err error
	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	var lapak Lapak
	if err := json.NewDecoder(r.Body).Decode(&lapak);err!=nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
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

	tx := db.DB.MustBegin()
	tx.MustExec("INSERT INTO stores (store_name,address,handphone,bank_number,id_user) VALUES (?, ? ,? , ?, ?)",lapak.NamaLapak,lapak.Alamat,lapak.Telepon,lapak.NomerRekening,id_user)
	var id int
	tx.Get(&id, "SELECT LAST_INSERT_ID() as id")
	jwt := jwt.MapClaims{
		"store_id": id,
		"store_name": lapak.NamaLapak,
	}
	if lapak.ApiKey, err = utils.GenerateToken(w, r, jwt, "storesecr3t"); err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}
	tx.MustExec("UPDATE stores SET api_key=? WHERE id=?", lapak.ApiKey, id)
	if err = tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error creating new lapak",http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, lapak, http.StatusOK, "success")
}



type Lapak struct {
	NamaLapak string `json:"namaLapak"`
	Alamat string `json:"alamat"`
	Telepon string `json:"telepon"`
	NomerRekening string `json:"nomerRekening"`
	ApiKey string `json:"apiKey"`
}
