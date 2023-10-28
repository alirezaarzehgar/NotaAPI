package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Asrez/NotaAPI/config"
)

var alertsDatabase map[string]string

func Alert(alertKey string) string {
	if alertsDatabase == nil {
		data, err := os.ReadFile(config.AlertDb())
		if err != nil {
			log.Println("read alert database: ", err)
			return ""
		}
		if err := json.Unmarshal(data, &alertsDatabase); err != nil {
			log.Println("unmarshal alerts: ", err)
			return ""
		}
	}
	return alertsDatabase[alertKey]
}
