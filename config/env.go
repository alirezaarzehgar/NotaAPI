package config

import (
	"os"
	"strconv"
)

func ListenerAddr() string {
	return os.Getenv("RUNNING_ADDR")
}

func JwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func LogDirectory() string {
	return os.Getenv("LOG_DIRECTORY")
}

type DbConf struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     uint64
}

func Db() (*DbConf, error) {
	port, err := strconv.ParseUint(os.Getenv("MYSQL_PORT"), 10, 64)
	if err != nil {
		return nil, err
	}

	conf := DbConf{
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DbName:   os.Getenv("MYSQL_DATABASE"),
		Port:     port,
	}
	return &conf, nil
}

func AlertDb() string {
	return os.Getenv("ALERT_DATABASE")
}

type AdminConfig struct {
	Username, Email, Password string
}

func Admin() AdminConfig {
	return AdminConfig{
		Username: os.Getenv("ADMIN_NAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASSWORD"),
	}
}

func Assets() string {
	return os.Getenv("ASSETS_DIRECTORY")
}
