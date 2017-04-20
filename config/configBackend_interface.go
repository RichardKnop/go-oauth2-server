package config

type ConfigBackend interface {
	LoadConfig() (*Config, error)
	RefreshConfig(newCnf *Config)
	InitConfigBackend()
}
