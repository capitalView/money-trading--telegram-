package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

var DbUser = goDotEnvVariable("DB_USER")
var DbPassword = goDotEnvVariable("DB_PASSWORD")
var DbName = goDotEnvVariable("DB_NAME")
var TokenTelegram = goDotEnvVariable("TELEGRAM_TOKEN")
var DbHost = goDotEnvVariable("DB_HOST")
var ChatIdAmin = goDotEnvVariable("CHAT_ID")
var ApiUrl = goDotEnvVariable("API_URL")
