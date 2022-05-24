package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err.Error())
	}

	return os.Getenv(key)
}
