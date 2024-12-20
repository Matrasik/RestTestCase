package handlers

import (
	"Test_REST/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"log"
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
	case "WITHDRAW":
		db.Withdraw(w, r, &wallet)
	}
	_, err = fmt.Fprint(w, wallet)
	if err != nil {
		return
	}
}

func (db *DBWrapper) CreateWallet(w http.ResponseWriter, r *http.Request, wallet *models.Wallet, tx *gorm.DB) error {
	wallet.Balance = 0
	err := tx.Exec(`
			INSERT INTO wallets (wallet_id, balance)
			VALUES ($1, $2)
			ON CONFLICT (wallet_id) DO NOTHING
		`, wallet.WalletId, wallet.Balance).Error
	if err != nil {
		tx.Rollback()
		log.Println("Error creating wallet:", err)
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (db *DBWrapper) Deposit(w http.ResponseWriter, r *http.Request, wallet *models.Wallet) {

	tx := db.DB.Begin()
	if tx.Error != nil {
		log.Println("Error begin transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var existingWallet models.Wallet
	err := tx.First(&existingWallet, "wallet_id = ?", wallet.WalletId).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		err := db.CreateWallet(w, r, wallet, tx)
		if err != nil {
			log.Println("Failed to create wallet", err)
		}
	} else if err != nil {
		tx.Rollback()
		log.Println("Error querying database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result := tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE wallet_id = $2", wallet.Amount, wallet.WalletId)
	if result.Error != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (db *DBWrapper) Withdraw(w http.ResponseWriter, r *http.Request, wallet *models.Wallet) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		log.Println("Error begin transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var currBalance int
	err := tx.Raw("SELECT balance FROM wallets WHERE wallet_id = $1", wallet.WalletId).Scan(&currBalance).Error
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if currBalance < wallet.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}
	result := db.DB.Exec("UPDATE wallets SET balance = balance - $1 WHERE wallet_id = $2",
		wallet.Amount, wallet.WalletId)
	if result.Error != nil {
		log.Printf("Error exec to database %s", tx.Error)
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx.Commit()

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
