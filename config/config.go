package config

import (
	"os"
	"strconv"
)

type Config struct {
	Email    EmailConfig
	Postgres PostgresConfig
}

type EmailConfig struct {
	SMTPServer  string
	SMTPPort    int
	Username    string
	Password    string
	FromAddress string
	FromName    string
	TLS         bool
}

type PostgresConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
}

func ReadConfigFromEnv() Config {
	return Config{

		Email: EmailConfig{
			SMTPServer:  os.Getenv("EMAIL_SMTP_SERVER"),
			SMTPPort:    getEnvAsInt("EMAIL_SMTP_PORT", 587),
			Username:    os.Getenv("EMAIL_USERNAME"),
			Password:    os.Getenv("EMAIL_PASSWORD"),
			FromAddress: os.Getenv("EMAIL_FROM_ADDRESS"),
			FromName:    os.Getenv("EMAIL_FROM_NAME"),
			TLS:         getEnvAsBool("EMAIL_TLS", true),
		},

		Postgres: PostgresConfig{
			Host:         os.Getenv("POSTGRES_HOST"),
			Port:         getEnvAsInt("POSTGRES_PORT", 5432),
			User:         os.Getenv("POSTGRES_USER"),
			Password:     os.Getenv("POSTGRES_PASSWORD"),
			DBName:       os.Getenv("POSTGRES_DBNAME"),
			SSLMode:      os.Getenv("POSTGRES_SSLMODE"),
			MaxIdleConns: getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 100),
		},
	}
}

// Helper function to get environment variable as integer
func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as boolean
func getEnvAsBool(name string, defaultValue bool) bool {
	valueStr := os.Getenv(name)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
