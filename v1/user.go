package v1

import (
	"net/http"
	"fmt"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
)

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ApiKey string `json:"apiKey"`
}

func CreateUser(w http.ResponseWriter, r *http.Request){
	var user User

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user);err != nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost);if err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("got error %s",err.Error()),http.StatusBadRequest)
	}

	user.Password = string(hashedPassword)



}