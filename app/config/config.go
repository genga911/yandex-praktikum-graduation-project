package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress          string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DataBaseURI         string `env:"DATABASE_URI" envDefault:""`
	ActualSystemAddress string `env:"ACTUAL_SYSTEM_ADDRESS" envDefault:""`
	JWTKey              string `env:"JWT_KEY" envDefault:"XVKjs6qaK9WiEr5g"`
}

func Get() *Config {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		panic(fmt.Sprintf("Parse config error: %s", err))
	}

	return &cfg
}

func InitFlags(cfg *Config) {
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "адрес и порт запуска сервиса")
	flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "адрес подключения к базе данных")
	flag.StringVar(&cfg.ActualSystemAddress, "f", cfg.ActualSystemAddress, "адрес системы расчёта начислений: переменная окружения ОС")
	flag.Parse()
}
