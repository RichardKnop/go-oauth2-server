package config

import "github.com/spf13/viper"

// NewConfig loads configuration from environment variables and returns *Config struct
func NewConfig() *Config {
	// Bind environment variables
	viper.BindEnv("oauth2_database_type")
	viper.BindEnv("oauth2_database_host")
	viper.BindEnv("oauth2_database_port")
	viper.BindEnv("oauth2_database_user")
	viper.BindEnv("oauth2_database_password")
	viper.BindEnv("oauth2_database_name")
	viper.BindEnv("oauth2_access_token_lifetime")
	viper.BindEnv("oauth2_refresh_token_lifetime")

	// Set sensible defaults for environment variables
	viper.SetDefault("oauth2_database_type", "postgres")
	viper.SetDefault("oauth2_database_host", "127.0.0.1")
	viper.SetDefault("oauth2_database_port", 5432)
	viper.SetDefault("oauth2_database_user", "go_oauth2_server")
	viper.SetDefault("oauth2_database_password", "")
	viper.SetDefault("oauth2_database_name", "go_oauth2_server")
	viper.SetDefault("oauth2_access_token_lifetime", 3600)     // 1 hour
	viper.SetDefault("oauth2_refresh_token_lifetime", 1209600) // 14 days

	return &Config{
		Database: DatabaseConfig{
			Type:         viper.GetString("oauth2_database_type"),
			Host:         viper.GetString("oauth2_database_host"),
			Port:         viper.GetInt("oauth2_database_port"),
			User:         viper.GetString("oauth2_database_user"),
			Password:     viper.GetString("oauth2_database_password"),
			DatabaseName: viper.GetString("oauth2_database_name"),
		},
		AccessTokenLifetime:  viper.GetInt("oauth2_access_token_lifetime"),
		RefreshTokenLifetime: viper.GetInt("oauth2_refresh_token_lifetime"),
	}
}
