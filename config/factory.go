package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/pkg/transport"
	"golang.org/x/net/context"
)

var (
	etcdEndpoints                         = "http://localhost:2379"
	etcdCertFile, etcdKeyFile, etcdCaFile string
	etcdConfigPath                        = "/config/go_oauth2_server.json"
	configLoaded                          bool
	dialTimeout                           = 5 * time.Second
	contextTimeout                        = 5 * time.Second
	reloadDelay                           = time.Second * 10
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

func init() {
	// Overwrite default values with environment variables if they are set
	if os.Getenv("ETCD_ENDPOINTS") != "" {
		etcdEndpoints = os.Getenv("ETCD_ENDPOINTS")
	}
	if os.Getenv("ETCD_CERT_FILE") != "" {
		etcdCertFile = os.Getenv("ETCD_CERT_FILE")
	}
	if os.Getenv("ETCD_KEY_FILE") != "" {
		etcdKeyFile = os.Getenv("ETCD_KEY_FILE")
	}
	if os.Getenv("ETCD_CA_FILE") != "" {
		etcdCaFile = os.Getenv("ETCD_CA_FILE")
	}
	if os.Getenv("ETCD_CONFIG_PATH") != "" {
		etcdConfigPath = os.Getenv("ETCD_CONFIG_PATH")
	}
}

// NewConfig loads configuration from etcd and returns *Config struct
// It also starts a goroutine in the background to keep config up-to-date
func NewConfig(mustLoadOnce bool, keepReloading bool) *Config {
	if configLoaded {
		return Cnf
	}

	// If the config must be loaded once successfully
	if mustLoadOnce && !configLoaded {
		// Read from remote config the first time
		newCnf, err := LoadConfig()
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
				<-time.After(reloadDelay)

				// Attempt to reload the config
				newCnf, err := LoadConfig()
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
func LoadConfig() (*Config, error) {
	cli, err := newEtcdClient(etcdEndpoints, etcdCertFile, etcdKeyFile, etcdCaFile)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// Read from remote config the first time
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	resp, err := cli.Get(ctx, etcdConfigPath)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			return nil, fmt.Errorf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			return nil, fmt.Errorf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			return nil, fmt.Errorf("client-side error: %v", err)
		default:
			return nil, fmt.Errorf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("key not found: %s", etcdConfigPath)
	}

	// Unmarshal the config JSON into the cnf object
	newCnf := new(Config)

	if err := json.Unmarshal([]byte(resp.Kvs[0].Value), newCnf); err != nil {
		return nil, err
	}

	return newCnf, nil
}

// RefreshConfig sets config through the pointer so config actually gets refreshed
func RefreshConfig(newCnf *Config) {
	*Cnf = *newCnf
}

func newEtcdClient(theEndpoints, certFile, keyFile, caFile string) (*clientv3.Client, error) {
	// Log the etcd endpoint for debugging purposes
	logger.Infof("ETCD Endpoints: %s", theEndpoints)

	// ETCD config
	etcdConfig := clientv3.Config{
		Endpoints:   strings.Split(theEndpoints, ","),
		DialTimeout: dialTimeout,
	}

	// Optionally, configure TLS transport
	if certFile != "" && keyFile != "" && caFile != "" {
		// Load client cert
		tlsInfo := transport.TLSInfo{
			CertFile:      certFile,
			KeyFile:       keyFile,
			TrustedCAFile: caFile,
		}

		// Setup HTTPS client
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		// Add TLS config
		etcdConfig.TLS = tlsConfig
	}

	// ETCD client
	return clientv3.New(etcdConfig)
}
