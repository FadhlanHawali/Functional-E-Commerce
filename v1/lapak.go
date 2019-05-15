package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
)

func CreateLapak(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}
}

type Lapak struct {
	ApiKey string `json:"apiKey"`
	Owner string `json:"owner"`
	NamaLapak string `json:"namaLapak"`
	Alamat string `json:"alamat"`
	Kontak Kontak `json:"kontak"`
}

type Kontak struct {
	Email string `json:"email"`
	Telepon int32 `json:"telepon"`
	NomerRekening int32 `json:"nomerRekening"`
}
