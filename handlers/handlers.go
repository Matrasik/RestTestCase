package handlers

import (
	"Test_REST/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type DBWrapper struct {
	DB *gorm.DB
}

func (db *DBWrapper) HandleDepWithdraw(w http.ResponseWriter, r *http.Request) {
	// Получить тип операции и вызывать функции
	ans, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	var wallet models.Wallet
	err = json.Unmarshal(ans, &wallet)
	if err != nil {
		return
	}
	switch wallet.OperationType {
	case "DEPOSIT":
		db.Deposit(w, r, &wallet)
		//TODO: вызов функции депозита
	case "WITHDRAW":
		db.Withdraw(w, r, &wallet)
		//TODO: вызов функции вноса
	}
	_, err = fmt.Fprint(w, wallet)
	if err != nil {
		return
	}
}

func (db *DBWrapper) Deposit(w http.ResponseWriter, r *http.Request, wallet *models.Wallet) {
	db.DB.Exec("UPDATE wallets SET balance = balance + $1 WHERE wallet_id = $2",
		wallet.Amount, wallet.WalletId)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (db *DBWrapper) Withdraw(w http.ResponseWriter, r *http.Request, wallet *models.Wallet) {
	db.DB.Exec("UPDATE wallets SET balance = balance - $1 WHERE wallet_id = $2",
		wallet.Amount, wallet.WalletId)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (db *DBWrapper) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//Получение баланса
	vars := mux.Vars(r)
	fmt.Fprint(w, vars["WALLET_UUID"])

	//TODO: обращение к бд и получение баланса по юид
}
