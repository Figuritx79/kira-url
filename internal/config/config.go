package config

import (
	"fmt"
	"kira-url/internal/env"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	Host           string
	Password       string
	User           string
	Port           int
	Name           string
	Schema         string
	SSLMode        string
	ChannelBinding string
	AutoMigrate    bool
}

type ServerConfig struct {
	Port        int
	Env         string
	Domain      string
	CorsDomains []string
}

type LoggerConfig struct {
	LogLevel string
}

func New() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        env.GetEnvInt("HTTP_PORT", 8080),
			Env:         env.GetEnvString("SERVER_ENV", "DEV"),
			Domain:      env.GetEnvString("SERVER_DOMAIN", ""),
			CorsDomains: env.GetEnvStringSlice("CORS_ORIGIN", []string{}),
		},
		Database: DatabaseConfig{
			Name:           env.GetEnvString("DB_DATABASE", "example"),
			Password:       env.GetEnvString("DB_PASSWORD", "your_password"),
			User:           env.GetEnvString("DB_USERNAME", "your_username"),
			Port:           env.GetEnvInt("DB_PORT", 3536),
			Host:           env.GetEnvString("DB_HOST", "localhost"),
			Schema:         env.GetEnvString("DB_SCHEMA", "public"),
			SSLMode:        env.GetEnvString("SSL_MODE", "require"),
			ChannelBinding: env.GetEnvString("CHANNEL_BINDING", "require"),
		},
		Logger: LoggerConfig{
			LogLevel: env.GetEnvString("LOG_LEVEL", "info"),
		},
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}
func (c *Config) validate() error {
	if c.Server.Port == 0 {
		return fmt.Errorf("HTTP_PORT is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Server.Domain == "" {
		return fmt.Errorf("SERVER_DOMAIN is required")
	}
	return nil
}
