package config

// DatabaseConfig stores database connection options
type DatabaseConfig struct {
	Type         string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string `json:"database_name"`
}

// Config stores configuration options
type Config struct {
	Database             DatabaseConfig
	AccessTokenLifetime  int `json:"access_token_lifetime"`
	RefreshTokenLifetime int `json:"refresh_token_lifetime"`
}
