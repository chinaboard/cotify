package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s environment variable is not set", key)
	}
	return val
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	return &Config{
		DBHost:     getEnvOrFatal("DB_HOST"),
		DBPort:     getEnvOrFatal("DB_PORT"),
		DBUser:     getEnvOrFatal("DB_USER"),
		DBPassword: getEnvOrFatal("DB_PASSWORD"),
		DBName:     getEnvOrFatal("DB_NAME"),
		ServerPort: getEnvOrFatal("SERVER_PORT"),
	}
}

func NewConfig(host, port, user, password, dbname, serverPort string) *Config {
	return &Config{
		DBHost:     host,
		DBPort:     port,
		DBUser:     user,
		DBPassword: password,
		DBName:     dbname,
		ServerPort: serverPort,
	}
}

// GetDSN assembles the DSN string from Config
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
