package utils

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

func GenerateToken(w http.ResponseWriter,r *http.Request, jwtMap jwt.MapClaims, secret string) (string,error){
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256,jwtMap)
	token, err := sign.SignedString([]byte(secret))
	if err != nil {
		return "",err
	}
	return token,nil
}

func ValidateToken(header string, secret string) (jwt.MapClaims,error){

	tokenString := header
	var claims jwt.MapClaims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if token != nil && err == nil {
		//TODO APA AJA YANG MAU DI CLAIM DARI TOKEN NYA
		//fmt.Println("token verified")
		claims = token.Claims.(jwt.MapClaims)
		//fmt.Println(claims)
		//mapstructure.Decode(claims["id"], &idDepartemen)
		//return idDepartemen,nil
	} else {
		return nil,err
	}

	return claims,nil

}