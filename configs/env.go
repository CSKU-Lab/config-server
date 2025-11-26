package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env map[string]string

func NewEnv() *env {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error loading .env file:", err)
	}

	return &env{
		"MONGO_URI":     os.Getenv("MONGO_URI"),
		"PORT":          os.Getenv("PORT"),
		"DATABASE_NAME": os.Getenv("DATABASE_NAME"),
	}
}

func (m *env) Get(key string) string {
	val, exists := (*m)[key]
	if !exists {
		log.Fatalf("Environment variable %s not found!", key)
	}
	return val
}
