package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
)


func (db *InDB) CreateLapak(w http.ResponseWriter, r *http.Request){

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

}
type Lapak struct {
	NamaLapak string `json:"namaLapak"`
	Alamat string `json:"alamat"`
	Telepon string `json:"telepon"`
	NomerRekening string `json:"nomerRekening"`
}
