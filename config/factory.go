package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	// Enable the remote features
	_ "github.com/spf13/viper/remote"
)

var configLoaded bool

// Let's start with some sensible defaults
var cnf = &Config{
	Database: DatabaseConfig{
		Type:         "postgres",
		Host:         "localhost",
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
		AuthCodeLifetime:     3600,    // TODO - should this be less than 1 hour?
	},
	Session: SessionConfig{
		Secret:   "test_secret",
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HTTPOnly: true,
	},
	TrustedClient: TrustedClientConfig{
		ClientID: "test_client",
		Secret:   "test_secret",
	},
}

// NewConfig loads configuration from etcd and returns *Config struct
// It also starts a goroutine in the background to keep config up-to-date
func NewConfig() *Config {
	if configLoaded {
		return cnf
	}

	// Bind etcd env vars
	viper.SetDefault("etcd_host", "localhost")
	viper.SetDefault("etcd_port", "2379")
	viper.BindEnv("etcd_host")
	viper.BindEnv("etcd_port")

	// Construct the ETCD URL
	etcdURL := fmt.Sprintf(
		"http://%s:%s",
		viper.Get("etcd_host"),
		viper.Get("etcd_port"),
	)

	// Config path
	configPath := "/config/go_oauth2_server.json"

	// Add a new ETCD remote provider
	runtimeViper := viper.New()
	runtimeViper.AddRemoteProvider("etcd", etcdURL, configPath)
	// Because there is no file extension in a stream of bytes
	runtimeViper.SetConfigType("json")

	// Read from remote config the first time.
	if err := runtimeViper.ReadRemoteConfig(); err != nil {
		log.Printf(
			"Unable to read remote config for the first time from %s/v2/keys%s: %v",
			etcdURL,
			configPath,
			err,
		)
	} else {
		runtimeViper.Unmarshal(&cnf)
	}

	// Open a goroutine to watch remote changes forever
	go func() {
		for {
			// Delay after each request
			time.Sleep(time.Second * 5)

			if err := runtimeViper.WatchRemoteConfig(); err != nil {
				log.Printf(
					"Unable to read remote config from %s/v2/keys%s: %v",
					etcdURL,
					configPath,
					err,
				)
				continue
			}

			// Unmarshal config
			runtimeViper.Unmarshal(&cnf)
		}
	}()

	configLoaded = true
	return cnf
}
