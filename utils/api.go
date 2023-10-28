package utils

import (
	"crypto/sha256"
	"encoding/hex"
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

func HashPassword(pass string) string {
	hashByte := sha256.Sum256([]byte(pass))
	hashStr := hex.EncodeToString(hashByte[:])
	return hashStr
}
