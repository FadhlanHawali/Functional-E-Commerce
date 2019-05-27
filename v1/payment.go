package v1

import (
	"net/http"
	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"encoding/json"
	"fmt"
	"bytes"
	"log"
	"github.com/spf13/viper"
	"github.com/gorilla/mux"
	"errors"
)

func CreatePayment (w http.ResponseWriter, r *http.Request, order OrderRepo, db *InDB) (interface{},error){

	if r.Method != "POST"{
		return nil,errors.New(http.StatusText(http.StatusMethodNotAllowed))
	}

	reqPayment := RequestPayment{
		WalletToken:fmt.Sprintf("%s",viper.Get("payment.walletToken")),
		Price:order.Total,
	}

	paymentBytes, err := json.Marshal(reqPayment);if err != nil {
		return nil,err
	}

	result,err := requestPayment(paymentBytes); if err!=nil{
		return nil,err
	}

	tx := db.DB.MustBegin()
	//TODO id customer
	tx.MustExec(fmt.Sprintf("UPDATE orders SET token_payment='%s' WHERE id_customer=%d",result.Data.Kode,order.Id_Customer))

	return map[string]interface{}{
		"url":issuePayment(result.Data.Kode),
	},nil
}

func requestPayment (bytesRepresentation []byte) (RespondPayment,error){
	var result RespondPayment
	req, err := http.NewRequest("POST", fmt.Sprintf("https://arta.ruangkarya.id/payment/create-bill"), bytes.NewBuffer(bytesRepresentation))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RespondPayment{},err
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("can't read response body : %s",err)
			return RespondPayment{},err
		} else {
			return result,err
		}
	}
}

func (db *InDB) UpdatePayment (w http.ResponseWriter, r *http.Request){

	if r.Method != "GET"{
		utils.WrapAPIError(w,r,http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
		return
	}
	token := mux.Vars(r)["user"]
	result,err := checkPayment(token);if err!=nil{
		utils.WrapAPIError(w,r,"error validating payment",http.StatusInternalServerError)
		return
	}

	data := result["data"].(map[string]interface{})
	isPaid := data["is_paid"].(bool)
	if isPaid == false{
		utils.WrapAPIError(w,r,"Bill is not paid yet",http.StatusOK)
		return
	}

	//status,err := updatePayment(db,token);if err!= nil{
	//
	//}
}

func issuePayment (token string) string{
	//TODO ADD SUCCESS REDIRECT
	return fmt.Sprintf("https://arta.ruangkarya.id/pay?paymentCode=%s",token)
}

func checkPayment (token string) (map[string]interface{},error) {
	var result map[string]interface{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://arta.ruangkarya.id/payment/get-status/%s",token), nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result,err
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("can't read response body : %s",err)
			return nil,err
		} else {
			return result,err
		}
	}
}

func updatePayment (db *InDB, token string) (bool,error){
	//tx := db.DB.MustBegin()

	return false,nil
}

type RequestPayment struct {
	WalletToken string `json:"token"`
	Price int `json:"jumlah"`
}

type RespondPayment struct {
	Status bool `json:"status"`
	Data DataPayment `json:"data"`
}

type DataPayment struct {
	Kode string `json:"kode"`
}