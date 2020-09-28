package main

import (
	"UserBalance/Services"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	db, err := sql.Open("mysql", "root:root@/avito_watch")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var balanceManager Services.BalanceManager
	balanceManager.DB = db



	var handler Services.Handler
	handler.BalanceManager = balanceManager

	router := mux.NewRouter()
	router.HandleFunc("/", handler.Index).Methods("GET")
	router.HandleFunc("/add", handler.HandleAdd).Methods("POST")
	router.HandleFunc("/users/{id}", handler.HandleUserInfo).Methods("GET")
	router.HandleFunc("/debit", handler.HandleDebit).Methods("POST")
	router.HandleFunc("/transfer", handler.HandleTransfer).Methods("POST")
	router.HandleFunc("/history/{id}", handler.HandleUserHistory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
