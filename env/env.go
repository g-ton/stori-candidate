package env

import (
	"os"

	"log"

	"github.com/joho/godotenv"
)

// Env config struct
type EnvApp struct {
	// Database Envs
	DB_HOST     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
	DB_SSLMODE  string
	DB_TIMEZONE string
}

// Get the env configuration
func GetEnv(env_file string) EnvApp {
	err := godotenv.Load(env_file)

	if err != nil {
		log.Printf("Error loading "+env_file+" file %v", err)
	}

	return EnvApp{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_SSLMODE:  os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE: os.Getenv("DB_TIMEZONE"),
	}
}
