package config

import (
	"log"
	"os"

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

	// Automatically read environment variables
	viper.AutomaticEnv()

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get the PostgreSQL DSN and replace environment variables
	postgresDSN := viper.GetString("database.dsn")
	postgresDSN = os.ExpandEnv(postgresDSN) // Ensure environment variables are expanded

	// Log the final DSN to debug the issue
	log.Println("Final Postgres DSN:", postgresDSN)

	// Return the config struct with updated values
	return Config{
		ServerPort:  viper.GetString("server.port"),
		PostgresDSN: postgresDSN,
	}
}
