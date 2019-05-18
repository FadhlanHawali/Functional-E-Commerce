package v1

import (
	"net/http"
	"fmt"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"log"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ApiKey string `json:"apiKey"`
}

func (db *InDB) CreateUser(w http.ResponseWriter, r *http.Request){
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
	log.Println(user.Password)
	jwt := jwt.MapClaims{
		"email":user.Email,
	}
	if user.ApiKey,err = utils.GenerateToken(w,r,jwt);err!=nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}
	tx := db.DB.MustBegin()

	if result := tx.MustExec(fmt.Sprintf("insert into users (email,password,api_key) values ('%s','%s','%s');",user.Email,user.Password,user.ApiKey));err!=nil{
		log.Println(fmt.Sprintf("got result %s ; error %s",result,err.Error()))
		utils.WrapAPIError(w,r,"error creating new user",http.StatusInternalServerError)
		return
	}else {
		tx.Commit()
	}
	tx.Commit()

	utils.WrapAPIError(w,r,"success creating new user",http.StatusCreated)
}