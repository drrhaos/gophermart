package configure

import (
	"flag"
	"gophermart/internal/logger"
	"net/url"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func (cfg *Config) ReadStartParams() bool {
	err := env.Parse(cfg)
	if err != nil {
		logger.Logger.Info("Не удалось найти переменные окружения ")
	}
	runAddress := flag.String("a", "127.0.0.1:8080", "адрес и порт запуска сервиса host:port")
	databaseURI := flag.String("d", "", "адрес подключения к базе данных postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable")
	accrualSystemAddress := flag.String("r", "", "адрес системы расчёта начислений")

	flag.Parse()
	if cfg.RunAddress == "" {
		cfg.RunAddress = *runAddress
	}
	if cfg.DatabaseURI == "" {
		cfg.DatabaseURI = *databaseURI
	}

	if cfg.AccrualSystemAddress == "" {
		cfg.AccrualSystemAddress = *accrualSystemAddress
	}

	_, errURL := url.ParseRequestURI("http://" + cfg.RunAddress)
	if errURL != nil {
		flag.PrintDefaults()
		return false
	}

	return true

}
