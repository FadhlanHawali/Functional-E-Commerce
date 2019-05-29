package v1

import (
	"net/http"
	"fmt"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"time"
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

	tx := db.DB.MustBegin()
	id := 0
	tx.Get(&id, "SELECT id FROM users WHERE email = ?", user.Email)
	if id > 0 {
		utils.WrapAPIError(w, r, "user already exist", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost);if err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("got error %s",err.Error()),http.StatusBadRequest)
	}
	user.Password = string(hashedPassword)
	tx.MustExec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
	tx.Get(&id, "SELECT LAST_INSERT_ID() as id")

	jwt := jwt.MapClaims{
		"id": id,
		"email":user.Email,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	}
	if user.ApiKey, err = utils.GenerateToken(w, r, jwt, "secret"); err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error creating new user", http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w , r, map[string]interface{} {
		"Email": user.Email,
		"ApiKey": user.ApiKey,
	}, http.StatusOK, "success")
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *InDB) Login(w http.ResponseWriter, r *http.Request){
	var user User

	if r.Method != "POST"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user);err != nil{
		utils.WrapAPIError(w,r,"Can't decode request body",http.StatusBadRequest)
		return
	}

	tx := db.DB.MustBegin()
	var hashedPassword string
	tx.Get(&hashedPassword, "SELECT password FROM users WHERE email = ?", user.Email)
	var id int
	tx.Get(&id, "SELECT id FROM users WHERE email = ?", user.Email)

	if (!CheckPasswordHash(user.Password, hashedPassword)) {
		utils.WrapAPIError(w, r, "wrong credential", http.StatusUnauthorized)
		return
	}

	jwt := jwt.MapClaims{
		"id": id,
		"email":user.Email,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	}

	key, err := utils.GenerateToken(w, r, jwt, "secret"); if err != nil{
		utils.WrapAPIError(w,r,fmt.Sprintf("error generating token. got error %s",err.Error()),http.StatusInternalServerError)
		return
	}

	user.ApiKey = key
	// query := `UPDATE users SET api_key=? WHERE id=?`
	// tx.MustExec(query, user.ApiKey, id)


	utils.WrapAPIData(w , r, map[string]interface{} {
		"Email": user.Email,
		"ApiKey": user.ApiKey,
	}, http.StatusOK, "success")
}