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
	SecretKey           string `env:"SECRET_KEY" envDefault:"XVKjs6qaK9WiEr5g"`
	CookieTTL           int    `env:"COOKIE_TTL" envDefault:"300"`
	AuthTTL             int    `env:"AUTH_TTL" envDefault:"86400"`
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

// GetAccuralRequestAddress number - идентификатор заказа
func (c *Config) GetAccuralRequestAddress(number string) string {
	return fmt.Sprintf("%s/api/orders/%s", c.ActualSystemAddress, number)
}
