package config

// DatabaseConfig stores database connection options
type DatabaseConfig struct {
	Type         string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

// Config stores configuration options
type Config struct {
	Database             DatabaseConfig
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	AuthCodeLifetime     int
}
