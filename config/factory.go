package config

import (
	"os"
	"time"

	"github.com/RichardKnop/go-oauth2-server/log"
)

var (
	configLoaded   bool
	dialTimeout    = 5 * time.Second
	contextTimeout = 5 * time.Second
	reloadDelay    = time.Second * 10
)

// Cnf ...
// Let's start with some sensible defaults
var Cnf = &Config{
	Database: DatabaseConfig{
		Type:         "postgres",
		Host:         "postgres",
		Port:         5432,
		User:         "go_oauth2_server",
		Password:     "",
		DatabaseName: "go_oauth2_server",
		MaxIdleConns: 5,
		MaxOpenConns: 5,
	},
	Oauth: OauthConfig{
		AccessTokenLifetime:  3600,    // 1 hour
		RefreshTokenLifetime: 1209600, // 14 days
		AuthCodeLifetime:     3600,    // 1 hour
	},
	Session: SessionConfig{
		Secret:   "test_secret",
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HTTPOnly: true,
	},
	IsDevelopment: true,
}

// NewConfig loads configuration from etcd and returns *Config struct
// It also starts a goroutine in the background to keep config up-to-date
func NewConfig(mustLoadOnce bool, keepReloading bool, backendType string) *Config {
	if configLoaded {
		return Cnf
	}

	var backend Backend

	switch backendType {
	case "etcd":
		backend = new(etcdBackend)
	case "consul":
		backend = new(consulBackend)
	default:
		log.FATAL.Printf("%s is not a valid backend", backendType)
		os.Exit(1)
	}

	backend.InitConfigBackend()

	// If the config must be loaded once successfully
	if mustLoadOnce && !configLoaded {
		// Read from remote config the first time
		newCnf, err := backend.LoadConfig()

		if err != nil {
			log.FATAL.Print(err)
			os.Exit(1)
		}

		// Refresh the config
		backend.RefreshConfig(newCnf)

		// Set configLoaded to true
		configLoaded = true
		log.INFO.Print("Successfully loaded config for the first time")
	}

	if keepReloading {
		// Open a goroutine to watch remote changes forever
		go func() {
			for {
				// Delay after each request
				<-time.After(reloadDelay)

				// Attempt to reload the config
				newCnf, err := backend.LoadConfig()
				if err != nil {
					log.ERROR.Print(err)
					continue
				}

				// Refresh the config
				backend.RefreshConfig(newCnf)

				// Set configLoaded to true
				configLoaded = true
				log.INFO.Print("Successfully reloaded config")
			}
		}()
	}

	return Cnf
}
