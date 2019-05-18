package utils

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
)

func WrapAPIError(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"Code":         code,
		"ErrorType":    http.StatusText(code),
		"ErrorDetails": message,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	}else {
		log.Println(fmt.Sprintf("can't wrap API error : %s",err))
	}
}

func WrapAPISuccess(w http.ResponseWriter, r *http.Request, message string, code int){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"Code":code,
		"Status":message,
	})
	if err==nil{
		log.Println(message)
		w.Write(result)
	}else{
		log.Println(fmt.Sprintf("can't wrap API success : %s",err))
	}
}