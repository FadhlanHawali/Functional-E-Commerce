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

	tx := db.DB.MustBegin()
	var id int
	tx.MustExec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
	tx.Get(&id, "SELECT LAST_INSERT_ID() as id")

	jwt := jwt.MapClaims{
		"id": id,
		"email":user.Email,
	}
	if user.ApiKey, err = utils.GenerateToken(w, r, jwt); err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}
	query := `UPDATE users SET api_key=? WHERE id=?`
	tx.MustExec(query, user.ApiKey, id)
	if err = tx.Commit(); err != nil {
		utils.WrapAPIError(w,r,"error creating new user",http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w,r,"success creating new user",http.StatusCreated)
}