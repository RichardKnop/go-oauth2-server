package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	etcdHost     = "localhost"
	etcdPort     = "2379"
	configPath   = "/config/go_oauth2_server.json"
	configLoaded bool
)

// Cnf ...
// Let's start with some sensible defaults
var Cnf = &Config{
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
func NewConfig(mustLoadOnce bool, keepReloading bool) *Config {
	if configLoaded {
		return Cnf
	}

	// Construct the ETCD endpoint
	etcdEndpoint := getEtcdEndpoint()
	logger.Infof("ETCD Endpoint: %s", etcdEndpoint)

	// ETCD config
	etcdClientConfig := client.Config{
		Endpoints: []string{etcdEndpoint},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	// ETCD client
	etcdClient, err := client.New(etcdClientConfig)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	// ETCD keys API
	kapi := client.NewKeysAPI(etcdClient)

	// If the config must be loaded once successfully
	if mustLoadOnce {
		// Read from remote config the first time
		newCnf, err := LoadConfig(kapi)
		if err != nil {
			logger.Fatal(err)
			os.Exit(1)
		}

		// Refresh the config
		RefreshConfig(newCnf)

		// Set configLoaded to true
		configLoaded = true
		logger.Info("Successfully loaded config for the first time")
	}

	if keepReloading {
		// Open a goroutine to watch remote changes forever
		go func() {
			for {
				// Delay after each request
				time.Sleep(time.Second * 10)

				// Attempt to reload the config
				newCnf, err := LoadConfig(kapi)
				if err != nil {
					logger.Error(err)
					continue
				}

				// Refresh the config
				RefreshConfig(newCnf)

				// Set configLoaded to true
				configLoaded = true
				logger.Info("Successfully reloaded config")
			}
		}()
	}

	return Cnf
}

// LoadConfig gets the JSON from ETCD and unmarshals it to the config object
func LoadConfig(kapi client.KeysAPI) (*Config, error) {
	// Read from remote config the first time
	resp, err := kapi.Get(context.Background(), configPath, nil)
	if err != nil {
		return nil, err
	}

	// Unmarshal the config JSON into the cnf object
	newCnf := new(Config)
	if err := json.Unmarshal([]byte(resp.Node.Value), newCnf); err != nil {
		return nil, err
	}

	return newCnf, nil
}

// RefreshConfig sets config through the pointer so config actually gets refreshed
func RefreshConfig(newCnf *Config) {
	*Cnf = *newCnf
}

// getEtcdURL builds ETCD endpoint from environment variables
func getEtcdEndpoint() string {
	// Construct the ETCD URL
	etcdHost := "localhost"
	if os.Getenv("ETCD_HOST") != "" {
		etcdHost = os.Getenv("ETCD_HOST")
	}
	etcdPort := "2379"
	if os.Getenv("ETCD_PORT") != "" {
		etcdPort = os.Getenv("ETCD_PORT")
	}
	return fmt.Sprintf("http://%s:%s", etcdHost, etcdPort)
}
