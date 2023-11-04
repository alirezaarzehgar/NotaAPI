package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"

	"github.com/Asrez/NotaAPI/config"
)

var writer *os.File

func setLogOutput() {
	var err error
	writer.Close()
	logpath := fmt.Sprintf("%s/%s.log", config.LogDirectory(), time.Now().Format("2006-01-02"))
	writer, err = os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Unable to open logfile:", err)
	}
	log.SetOutput(writer)
}

func InitLogger() {
	setLogOutput()
	log.SetFlags(log.Ltime | log.Lshortfile)
	gocron.Every(1).Day().Do(setLogOutput)
	gocron.Start()
}

func GetLogger() *log.Logger {
	return log.New(writer, "", log.Ltime|log.Lshortfile)
}

func stopLogger() {
	writer.Close()
}

func DebugLog(v ...any) {
	if config.Debug() {
		log.Println(v...)
	}
}
