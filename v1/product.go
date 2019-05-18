package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
)

func (db *InDB) AddProduct (w http.ResponseWriter, r *http.Request){

	var product []Product

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product);err!=nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

}

type Product struct {
	NamaBarang string `json:"namaBarang"`
	Deskripsi string `json:"deskripsi"`
	Quantity string `json:"quantity"`
	Harga int `json:"harga"`
	UrlGambar string `json:"urlGambar"`
}
