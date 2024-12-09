package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Host            string
	Port            string
	MetricsInterval int
}

func LoadConfig() *Config {
	host := os.Getenv("EXPORTER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("EXPORTER_PORT")
	if port == "" {
		port = "8080"
	}

	interval := 1
	intervalEnv := os.Getenv("EXPORTER_METRICS_INTERVAL")
	if intervalEnv != "" {
		if _, err := fmt.Sscanf(intervalEnv, "%d", &interval); err != nil {
			log.Printf("Неверный EXPORTER_METRICS_INTERVAL: %v. Используется значение по умолчанию: 10", err)
		}
	}

	return &Config{
		Host:            host,
		Port:            port,
		MetricsInterval: interval,
	}
}

func ValidateConfig(cfg *Config) {
	if cfg.Port == "" {
		log.Fatal("EXPORTER_PORT не может быть пустым")
	}
}
