package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(){
	err:= godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("environment variables not found")
	}
	return value
}