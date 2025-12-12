package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `env:"LEVEL"`
	HTTP         HTTP
	Storage      Storage
	URLShortener URLShortener
}

type HTTP struct {
	Host        string `env:"SERVER_HOST"`
	Addr        string `env:"SERVER_ADDR"`
	Timeout     string `env:"SERVER_TIMEOUT_SECONDS"`
	IdleTimeout string `env:"SERVER_IDLE_TIMEOUT_SECONDS"`
	MaxConn     string `env:"SERVER_MAX_CONN"`
}

type Storage struct {
	Type string `env:"DB_TYPE"`
	URL  string `env:"DB_URL"`
}

type URLShortener struct {
	TTL    time.Duration `env:"URL_TTL"`
	Length int           `env:"URL_LENGTH"`
	Salt   string        `env:"URL_SALT"`
}

func MustLoad(path string) *Config {
	var cfg Config

	_ = cleanenv.ReadConfig(path, &cfg)

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read environment variables: %v", err)
	}

	return &cfg
}
