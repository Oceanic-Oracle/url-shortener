package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	PATH = "../.env"
)

type Config struct {
	Env          string `env:"ENV"`
	HTTP         HTTP
	Storage      Storage
	URLShortener URLShortener
}

type HTTP struct {
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

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(PATH, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	return &cfg
}
