package config

import "fmt"

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Telegram TelegramConfig `toml:"telegram"`
	Postgres PostgresConfig `toml:"postgres"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	WebhookURL string `toml:"webhook_url"`
	Port       int    `toml:"port"`
}

// TelegramConfig holds Telegram bot configuration
type TelegramConfig struct {
	BotToken string `toml:"bot_token"`
}

// PostgresConfig holds database connection configuration
type PostgresConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	SSLMode  string `toml:"sslmode"`
}

// DSN returns the PostgreSQL connection string
func (pc PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pc.Host, pc.Port, pc.User, pc.Password, pc.DBName, pc.SSLMode)
}
