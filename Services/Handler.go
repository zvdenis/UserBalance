package Services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const  (
	priceAscendingSort = "priceAscending"
	priceDescendingSort = "priceDescending"
)

type RequestData struct {
	UserID       int    `json:"userID"`
	SecondUserID int    `json:"secondUserID"`
	Value        int    `json:"value"` // Для валют предпочтителен Int. Если необходимо, можно добавить второй Int для копеек
	Message      string `json:message`
}

type Response struct {
	Message string `json:"message"`
}

//Отвечает за оброаботку запроосов
type Handler struct {
	BalanceManager BalanceManager
}

//Базовая страница
func (handler Handler) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

//Обработка добваления денег
func (handler Handler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data RequestData
	data.Message = "not provided"
	_ = json.NewDecoder(r.Body).Decode(&data)
	if data.Value <= 0 || data.UserID <= 0 {
		json.NewEncoder(w).Encode(Response{Message: "Invalid format"})
		return
	}

	err := handler.BalanceManager.AddMoney(data.UserID, data.Value, data.Message)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "Failed"})
		return
	}
	json.NewEncoder(w).Encode(Response{Message: "Success"})
}

//Обработка информации о пользователе
func (handler Handler) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	userID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "incorrect id"})
		return
	}
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "RUB"
	}
	if err != nil {
		json.NewEncoder(w).Encode(Response{Message: "Invalid currency key"})
		return
	}
	money, err := handler.BalanceManager.GetConvertedUserMoney(userID, currency)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "Unable to get balance"})
		return
	}

	json.NewEncoder(w).Encode(Response{Message: money})
}

//Обработка снятия денег
func (handler Handler) HandleDebit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data RequestData
	data.Message = "not provided"
	_ = json.NewDecoder(r.Body).Decode(&data)
	if data.Value <= 0 || data.UserID <= 0 {
		json.NewEncoder(w).Encode(Response{Message: "Invalid format"})
		return
	}

	err := handler.BalanceManager.DebitMoney(data.UserID, data.Value, data.Message)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "Unable to debit"})
		return
	}
	json.NewEncoder(w).Encode(Response{Message: "Success"})
}

//Обработка перевода денег
func (handler Handler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data RequestData
	_ = json.NewDecoder(r.Body).Decode(&data)
	if data.Value <= 0 || data.UserID <= 0 || data.SecondUserID <= 0 {
		json.NewEncoder(w).Encode(Response{Message: "Invalid format"})
		return
	}

	err := handler.BalanceManager.TransferMoney(data.UserID, data.SecondUserID, data.Value, data.Message)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "Unable to debit"})
		return
	}
	json.NewEncoder(w).Encode(Response{Message: "Success"})
}

//Обработка истории переводов пользователя
func (handler Handler) HandleUserHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	userID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "incorrect id"})
		return
	}

	list, err := handler.BalanceManager.GetUserHistory(userID)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{Message: "Unable to get balance"})
		return
	}

	json.NewEncoder(w).Encode(list)
}
