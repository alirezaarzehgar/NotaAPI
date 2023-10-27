package main

import (
	"log"

	"github.com/Asrez/NotaAPI/config"
	"github.com/Asrez/NotaAPI/routes"
	"github.com/Asrez/NotaAPI/utils"
	"github.com/joho/godotenv"
)

func main() {
	go utils.HandleSignalInterrupt()
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("faild to load .env: ", err)
	}
	utils.InitLogger()
	log.Println("Start application")
	if err := routes.Init().Start(config.ListenerAddr()); err != nil {
		log.Print("echo start:", err)
	}
}
