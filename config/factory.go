package config

import "github.com/spf13/viper"

// NewConfig loads configuration from environment variables and returns *Config struct
func NewConfig() *Config {
	// Bind environment variables
	viper.BindEnv("database_type")
	viper.BindEnv("database_host")
	viper.BindEnv("database_port")
	viper.BindEnv("database_user")
	viper.BindEnv("database_password")
	viper.BindEnv("database_name")
	viper.BindEnv("access_token_lifetime")
	viper.BindEnv("refresh_token_lifetime")

	// Set sensible defaults for environment variables
	viper.SetDefault("database_type", "postgres")
	viper.SetDefault("database_host", "127.0.0.1")
	viper.SetDefault("database_port", 5432)
	viper.SetDefault("database_user", "go_microservice_example")
	viper.SetDefault("database_password", "")
	viper.SetDefault("database_name", "go_microservice_example")
	viper.SetDefault("access_token_lifetime", 3600)     // 1 hour
	viper.SetDefault("refresh_token_lifetime", 1209600) // 14 days

	return &Config{
		Database: DatabaseConfig{
			Type:         viper.GetString("database_type"),
			Host:         viper.GetString("database_host"),
			Port:         viper.GetInt("database_port"),
			User:         viper.GetString("database_user"),
			Password:     viper.GetString("database_password"),
			DatabaseName: viper.GetString("database_name"),
		},
		AccessTokenLifetime:  viper.GetInt("access_token_lifetime"),
		RefreshTokenLifetime: viper.GetInt("refresh_token_lifetime"),
	}
}
