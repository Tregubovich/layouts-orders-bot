package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	TG struct {
		Host       string `env:"TG_HOST" envDefault:"api.telegram.org"`
		TgBotToken string `env:"TG_BOT_TOKEN"`
	}

	SQLite struct {
		Path string `env:"SQLITE_PATH" envDefault:"./data/sqlite/storage.db"`
	}

	BatchSize int `env:"BATCH_SIZE" envDefault:"100"`
}

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
