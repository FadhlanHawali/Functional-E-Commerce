package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
)

const TokenContextKey = "MyAppToken"

func (db *InDB) CreateLapak(w http.ResponseWriter, r *http.Request){

	var lapak Lapak
	var id int
	var err error
	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}


	if err := json.NewDecoder(r.Body).Decode(&lapak);err!=nil{
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
	tx.MustExec("INSERT INTO stores (store_name,address,handphone,bank_number,id_user) VALUES (?, ? ,? , ?, ?)",lapak.NamaLapak,lapak.Alamat,lapak.Telepon,lapak.NomerRekening,id)

	if err = tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error creating new lapak",http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w,r,"success creating new lapak",http.StatusCreated)
}



type Lapak struct {
	NamaLapak string `json:"namaLapak"`
	Alamat string `json:"alamat"`
	Telepon string `json:"telepon"`
	NomerRekening string `json:"nomerRekening"`
}
