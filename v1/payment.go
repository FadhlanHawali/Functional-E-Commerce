package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/FadhlanHawali/Functional-E-Commerce/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func CreatePayment(w http.ResponseWriter, r *http.Request, order OrderRepo, db *InDB) (interface{}, error) {

	if r.Method != "POST" {
		return nil, errors.New(http.StatusText(http.StatusMethodNotAllowed))
	}

	reqPayment := RequestPayment{
		WalletToken: fmt.Sprintf("%s", viper.Get("payment.walletToken")),
		Price:       order.Total,
	}

	paymentBytes, err := json.Marshal(reqPayment)
	if err != nil {
		return nil, err
	}

	result, err := requestPayment(paymentBytes)
	if err != nil {
		return nil, err
	}
	tx := db.DB.MustBegin()
	//TODO id customer
	log.Println("ID",order.Id)
	log.Println("ID Cust",order.Id_Customer)
	tx.MustExec(fmt.Sprintf("UPDATE orders SET token_payment='%s' WHERE id=%d and id_customer=%d", result.Data.Kode, order.Id,order.Id_Customer))
	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error creating new order", http.StatusInternalServerError)
		return nil,err
	}
	return map[string]interface{}{
		"url": issuePayment(result.Data.Kode, order.Id, order.Id_Customer),
	}, nil
}

func requestPayment(bytesRepresentation []byte) (RespondPayment, error) {
	var result RespondPayment
	req, err := http.NewRequest("POST", fmt.Sprintf("https://arta.ruangkarya.id/payment/create-bill"), bytes.NewBuffer(bytesRepresentation))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RespondPayment{}, err
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("can't read response body : %s", err)
			return RespondPayment{}, err
		} else {
			return result, err
		}
	}
}

func (db *InDB) UpdatePayment(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Println("Update Payment Jalan")
	token := mux.Vars(r)["token"]
	result, err := checkPayment(token)
	if err != nil {
		utils.WrapAPIError(w, r, "error validating payment", http.StatusInternalServerError)
		return
	}
	log.Println("RESULT : ",result)

	data := result["data"].(map[string]interface{})
	isPaid := data["is_paid"].(bool)
	if isPaid == false {
		utils.WrapAPIError(w, r, "Bill is not paid yet", http.StatusBadRequest)
		return
	}

	idCustomer := mux.Vars(r)["idCustomer"]
	idOrder := mux.Vars(r)["idOrder"]
	tx := db.DB.MustBegin()
	tx.MustExec("UPDATE orders SET status=2 WHERE id = ? and id_customer = ?", idOrder, idCustomer)

	if err := tx.Commit(); err != nil {
		utils.WrapAPIError(w, r, "error updating token", http.StatusInternalServerError)
		return
	}
	utils.WrapAPISuccess(w, r, "success do payment", http.StatusOK)
	return
	//status,err := updatePayment(db,token);if err!= nil{
	//
	//}
}

func issuePayment(token string, idOrder int, idCustomer int) string {
	//TODO ADD SUCCESS REDIRECT

	return fmt.Sprintf("https://arta.ruangkarya.id/pay?paymentCode=%s&successRedirect=http://%s:4321/v1/store/order/%d/user/%d/payment/%s", token, viper.Get("host.ip"), idOrder, idCustomer, token)
}

func checkPayment(token string) (map[string]interface{}, error) {
	var result map[string]interface{}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://arta.ruangkarya.id/payment/get-status/%s", token), nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("can't read response body : %s", err)
			return nil, err
		} else {
			return result, err
		}
	}
}

func updatePayment(db *InDB, token string) (bool, error) {
	//tx := db.DB.MustBegin()

	return false, nil
}

type RequestPayment struct {
	WalletToken string `json:"token"`
	Price       int    `json:"jumlah"`
}

type RespondPayment struct {
	Status bool        `json:"status"`
	Data   DataPayment `json:"data"`
}

type DataPayment struct {
	Kode string `json:"kode"`
}
