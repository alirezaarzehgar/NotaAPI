package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignalInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("Stop application")
	stopLogger()
	os.Exit(1)
}
