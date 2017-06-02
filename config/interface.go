package config

// Backend defines a configuration backend, implement this interface
// to support additional backends
type Backend interface {
	LoadConfig() (*Config, error)
	RefreshConfig(newCnf *Config)
	InitConfigBackend()
}
