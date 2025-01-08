package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string
	PostgresDSN string
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return Config{
		ServerPort:  viper.GetString("server.port"),
		PostgresDSN: viper.GetString("database.dsn"),
	}
}
