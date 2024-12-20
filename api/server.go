package api

import (
	"Test_REST/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

const PORT = ":8888"

func StartServer(db *gorm.DB) {
	dataBase := handlers.DBWrapper{DB: db}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/wallet", dataBase.HandleDepWithdraw).Methods("POST")
	router.HandleFunc("/api/v1/wallets/{WALLET_UUID}", dataBase.GetBalance).Methods("GET")
	server := &http.Server{
		Addr:         PORT,
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	err := http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		log.Fatal("Error starting server")
		return
	}
}
