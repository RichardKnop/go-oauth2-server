package config

import (
	"github.com/spf13/viper"
)

// Factory loads configuration from environment variables (with some
// sensible defaults for dev environment) and returns *Config struct
func Factory() *Config {
	// Set defaults for environment variables
	viper.SetDefault("DATABASE_TYPE", "postgres")
	viper.SetDefault("DATABASE_HOST", "127.0.0.1")
	viper.SetDefault("DATABASE_PORT", 5432)
	viper.SetDefault("DATABASE_USER", "richardknop")
	viper.SetDefault("DATABASE_PASSWORD", "")
	viper.SetDefault("DATABASE_NAME", "go_microservice_example")
	viper.SetDefault("ACCESS_TOKEN_LIFETIME", 3600)
	viper.SetDefault("REFRESH_TOKEN_LIFETIME", 3600)

	return &Config{
		Database: DatabaseConfig{
			Type:         viper.GetString("DATABASE_TYPE"),
			Host:         viper.GetString("DATABASE_HOST"),
			Port:         viper.GetInt("DATABASE_PORT"),
			User:         viper.GetString("DATABASE_USER"),
			Password:     viper.GetString("DATABASE_PASSWORD"),
			DatabaseName: viper.GetString("DATABASE_NAME"),
		},
		AccessTokenLifetime:  viper.GetInt("ACCESS_TOKEN_LIFETIME"),
		RefreshTokenLifetime: viper.GetInt("REFRESH_TOKEN_LIFETIME"),
	}
}
