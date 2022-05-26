package configuration

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DBDSN                  string `env:"DBDSN"`
	ACCESS_TOKEN_SECRET    string `env:"ACCESS_TOKEN_SECRET"`
	REFRESH_TOKEN_SECRET   string `env:"REFRESH_TOKEN_SECRET"`
	ACCESS_TOKEN_LIFETIME  int    `env:"ACCESS_TOKEN_LIFETIME"`
	REFRESH_TOKEN_LIFETIME int    `env:"ACCESS_TOKEN_LIFETIME"`
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var c Config
	err := env.Parse(&c)
	if err != nil {
		panic(err)
	}
	return c
}
