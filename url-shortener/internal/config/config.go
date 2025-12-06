package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	PATH = "../.env"
)

type Config struct {
	Env        string `env:"ENV"`
	Htppserver Httpserver
	Storage    Storage
}

type Httpserver struct {
	Addr         string `env:"SERVER_ADDR"`
	Timeout      string `env:"SERVER_TIMEOUT_SECONDS"`
	IddleTimeout string `env:"SERVER_IDLE_TIMEOUT_SECONDS"`
	MaxConn      string `env:"SERVER_MAX_CONN"`
}

type Storage struct {
	Type string `env:"DB_TYPE"`
	URL  string `env:"DB_URL"`
}

type URLShortener struct {
	TTL    int `env:"URL_TTL_MINUTES"`
	Length int `env:"URL_LENGTH"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(PATH, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	return &cfg
}
