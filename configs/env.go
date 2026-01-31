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
		"MONGO_URI":            os.Getenv("MONGO_URI"),
		"PORT":                 os.Getenv("PORT"),
		"DATABASE_NAME":        os.Getenv("DATABASE_NAME"),
		"TASK_SERVER_URL":      os.Getenv("TASK_SERVER_URL"),
		"GO_GRADER_SERVER_URL": os.Getenv("GO_GRADER_SERVER_URL"),
		"REDIS_SERVER_URL":     os.Getenv("REDIS_SERVER_URL"),
		"REDIS_PASSWORD":       os.Getenv("REDIS_PASSWORD"),
	}
}

func (m *env) Get(key string) string {
	val, exists := (*m)[key]
	if !exists {
		log.Fatalf("Environment variable %s not found!", key)
	}
	return val
}
