package main

import (
	"Test_REST/api"
	"Test_REST/db"
	"Test_REST/db/migration"
	"log"
)

func main() {
	dataBase, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed connect to database")
	}
	migration.Migrate(dataBase)
	api.StartServer(dataBase)
}
