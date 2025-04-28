// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost string
	ServerPort int
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	LogLevel   string
}

// LoadConfig читает конфигурацию из файла .env
func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // поддержка переменных окружения

	// Если файла нет, можно проигнорировать ошибку (тогда читаем только из ENV)
	_ = viper.ReadInConfig()

	// Установим значения по умолчанию
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("LOG_LEVEL", "info")

	cfg := &Config{
		ServerHost: viper.GetString("SERVER_HOST"),
		ServerPort: viper.GetInt("SERVER_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBName:     viper.GetString("DB_NAME"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		LogLevel:   viper.GetString("LOG_LEVEL"),
	}
	return cfg, nil
}
