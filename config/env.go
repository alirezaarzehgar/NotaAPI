package config

import "os"

func ListenerAddr() string {
	return os.Getenv("RUNNING_ADDR")
}

func JwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func LogDirectory() string {
	return os.Getenv("LOG_DIRECTORY")
}
