package config

import (
	"fmt"
	"log"
	"time"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Env  string `mapstructure:"APP_ENV"`
	Port string `mapstructure:"APP_PORT"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"JWT_SECRET"`
	Expire time.Duration `mapstructure:"JWT_EXPIRE"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetDefault("APP_NAME", "go-rbac")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "rbac_db")
	viper.SetDefault("JWT_EXPIRE", "24h")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading .env file: %v\n", err)
	}

	config := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
	}

	// Parse JWT expiration duration
	if expStr := viper.GetString("JWT_EXPIRE"); expStr != "" {
		duration, err := time.ParseDuration(expStr)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRE duration: %w", err)
		}
		config.JWT.Expire = duration
	}

	if config.Database.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if config.Database.Port == "" {
		return nil, fmt.Errorf("DB_PORT is required")
	}
	if config.Database.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if config.Database.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return config, nil
}

// GetDatabaseURL returns the formatted database connection string
func (c *DatabaseConfig) GetDatabaseURL() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
	)
	return dsn
}