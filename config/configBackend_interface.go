package config

// ConfigBackend defines the exported methods
type ConfigBackend interface {
	LoadConfig() (*Config, error)
	RefreshConfig(newCnf *Config)
	InitConfigBackend()
}
