package configuration

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	STORAGE_PATH string `env:"STORAGE_PATH"`
	GRPC_ADDR    string `env:"GRPC_ADDR"`
	GRPC_TIMEOUT int    `env:"GRPC_TIMEOUT"`
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
