package config

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"os"
)

func Get() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error when load env %s", err.Error())
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
		},
		Redis: Redis{
			Addr: os.Getenv("REDIS_HOST"),
			Pass: os.Getenv("REDIS_PASS"),
		},
		Email: Email{
			Host:     os.Getenv("EMAIL_HOST"),
			Port:     os.Getenv("EMAIL_PORT"),
			User:     os.Getenv("EMAIL_USER"),
			Password: os.Getenv("EMAIL_PASS"),
		},
	}
}
