package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func CreateLapak(w http.ResponseWriter, r *http.Request){
	var lapak Lapak

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&lapak);err != nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(lapak.Password), bcrypt.DefaultCost);if err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("got error %s",err.Error()),http.StatusBadRequest)
	}

	lapak.Password = string(hashedPassword)

}

type Lapak struct {
	ApiKey string `json:"apiKey"`
	Owner string `json:"owner"`
	NamaLapak string `json:"namaLapak"`
	Alamat string `json:"alamat"`
	Kontak Kontak `json:"kontak"`
	Password string `json:"password"`
}

type Kontak struct {
	Email string `json:"email"`
	Telepon int32 `json:"telepon"`
	NomerRekening int32 `json:"nomerRekening"`
}
