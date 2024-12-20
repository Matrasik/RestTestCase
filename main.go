package main

import (
	"Test_REST/api"
	"Test_REST/db"
	"Test_REST/db/migration"
	"log"
)

//func init() {
//	if err := godotenv.Load(); err != nil {
//		log.Print("No .env file found")
//	}
//}

func main() {
	dataBase, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed connect to database")
	}
	migration.Migrate(dataBase)
	api.StartServer(dataBase)
}
