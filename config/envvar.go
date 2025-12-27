package config

import (
	"os"
	"strconv"
)

type evBackend struct{}

func (b *evBackend) InitConfigBackend() {
}

// LoadConfig gets the JSON from ETCD and unmarshals it to the config object
func (b *evBackend) LoadConfig() (*Config, error) {

	newCnf := new(Config)
	newCnf.Database.Type = "postgres"
	newCnf.Database.Host = os.Getenv("PGHOST")
	newCnf.Database.User = os.Getenv("POSTGRES_USER")
	newCnf.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	newCnf.Database.DatabaseName = os.Getenv("POSTGRES_DB")

	port, portExists := os.LookupEnv("POSTGRES_PORT")
	if portExists {
		if n, err := strconv.Atoi(port); err == nil {
			newCnf.Database.Port = n
		} else {
			newCnf.Database.Port = 5432
		}
	} else {
		newCnf.Database.Port = 5432
	}

	newCnf.Database.MaxOpenConns = 5
	newCnf.Database.MaxIdleConns = 5

	newCnf.Oauth.AccessTokenLifetime = 3600
	newCnf.Oauth.RefreshTokenLifetime = 1209600
	newCnf.Oauth.AuthCodeLifetime = 3600

	secret, secretExists := os.LookupEnv("SESSION_SECRET")
	if secretExists {
		newCnf.Session.Secret = secret
	} else {
		newCnf.Session.Secret = "test_secret"
	}
	newCnf.Session.Path = "/"
	newCnf.Session.MaxAge = 604800
	newCnf.Session.HTTPOnly = true

	newCnf.IsDevelopment = true

	return newCnf, nil
}

// RefreshConfig sets config through the pointer so config actually gets refreshed
func (b *evBackend) RefreshConfig(newCnf *Config) {
	*Cnf = *newCnf
}
