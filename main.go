package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Asrez/NotaAPI/api/handlers"
	"github.com/Asrez/NotaAPI/api/routes"
	"github.com/Asrez/NotaAPI/config"
	"github.com/Asrez/NotaAPI/database"
	"github.com/Asrez/NotaAPI/utils"
)

func main() {
	go utils.HandleSignalInterrupt()
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("faild to load .env: ", err)
	}
	utils.InitLogger()

	dbConf, err := config.Db()
	if err != nil {
		log.Fatal(".env: ", err)
	}

	db, err := database.Init(dbConf)
	if err != nil {
		log.Fatal("database: ", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("migrate: ", err)
	}

	handlers.SetDB(db)

	log.Println("Start application")
	if err := routes.Init().Start(config.ListenerAddr()); err != nil {
		log.Print("echo start:", err)
	}
}
