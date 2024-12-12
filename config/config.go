package config

import (
	"os"
	"strconv"
	"time"
)

var (
	SECRET_KEY = os.Getenv("SECRET_KEY")

	SERVICE_HOST = os.Getenv("SERVICE_HOST")
	SERVICE_PORT = getInt32("SERVICE_PORT")
	
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = getInt32("DB_PORT")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")

	BCRYPT_COAST = 3

	JWT_DURATION = time.Hour * 12 // half of day

	MAIL_HOST = os.Getenv("MAIL_HOST")
	MAIL_PORT = getInt32("MAIL_PORT")
	MAIL_USER = os.Getenv("MAIL_USER")
	MAIL_PASSWORD = os.Getenv("MAIL_PASSWORD")
)

func getInt32(envName string) int32 {
	value, _ := strconv.Atoi(os.Getenv(envName))
	return int32(value)
}
