package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
)

const TokenContextKey = "MyAppToken"

func (db *InDB) CreateLapak(w http.ResponseWriter, r *http.Request){

	//var lapak Lapak
	//
	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	//if err := json.NewDecoder(r.Body).Decode(&lapak);err!=nil{
	//	utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
	//	return
	//}
	//
	//tx := db.DB.MustBegin()

	//if token := r.Context().Value(TokenContextKey); token != nil {
	//	// User is logged in
	//	log.Printf("TOKEN : %s",token)
	//} else {
	//	// User is not logged in
	//	log.Println("GAADA TOKEN")
	//}
}



type Lapak struct {
	NamaLapak string `json:"namaLapak"`
	Alamat string `json:"alamat"`
	Telepon string `json:"telepon"`
	NomerRekening string `json:"nomerRekening"`
}
