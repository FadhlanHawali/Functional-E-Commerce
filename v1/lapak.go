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
