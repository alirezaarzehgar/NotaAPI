package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Asrez/NotaAPI/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var alertsDatabase map[string]string

const DATE_FORMAT = "2006-01-02"

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

func ReturnAlert(c echo.Context, status int, alertKey string, extra ...any) error {
	DebugLog("response error:", status, "cause:", extra)
	return c.JSON(status, map[string]any{
		"status": false,
		"alert":  Alert(alertKey),
		"data":   []any{},
	})
}

func CreateSHA256(pass string) string {
	hashByte := sha256.Sum256([]byte(pass))
	hashStr := hex.EncodeToString(hashByte[:])
	return hashStr
}

var EXPTIME = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30))

func CreateUserToken(id uint, email, user string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprint(id),
		Issuer:    email,
		Subject:   user,
		ExpiresAt: EXPTIME,
	})
	bearer, _ := token.SignedString(config.JwtSecret())
	return bearer
}

func CreateGuestToken() string {
	rData := make([]byte, 10)
	if _, err := rand.Read(rData); err != nil {
		log.Println("rand.Read(): ", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   string(rData),
		ExpiresAt: EXPTIME,
	})
	bearer, _ := token.SignedString(config.JwtSecret())
	return bearer
}

func GetToken(c echo.Context) string {
	bearer := c.Request().Header.Get("Authorization")
	return bearer[len("Bearer "):]
}

func GetUserId(c echo.Context) uint {
	bearer := c.Request().Header.Get("Authorization")
	token, _, _ := new(jwt.Parser).ParseUnverified(bearer[len("Bearer "):], jwt.MapClaims{})
	claims := token.Claims.(jwt.MapClaims)

	_, ok := claims["jti"]
	if !ok {
		return 0
	}

	id, _ := strconv.Atoi(claims["jti"].(string))
	return uint(id)

}

func GetUnavailableStoriesFilter() {

}
