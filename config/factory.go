package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/pkg/transport"
)

var (
	etcdEndpoints                         = "http://localhost:2379"
	etcdCertFile, etcdKeyFile, etcdCaFile string
	etcdConfigPath                        = "/config/go_oauth2_server.json"
	configLoaded                          bool
)

// Cnf ...
// Let's start with some sensible defaults
var Cnf = &Config{
	Database: DatabaseConfig{
		Type:         "postgres",
		Host:         "localhost",
		Port:         5432,
		User:         "go_oauth2_server",
		Password:     "go_oauth2_server",
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

	// Init ETCD client
	etcdClient, err := newEtcdClient(etcdEndpoints, etcdCertFile, etcdKeyFile, etcdCaFile)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	// ETCD keys API
	//kapi := clientv3.NewKeysAPI(*etcdClient)

	// If the config must be loaded once successfully
	if mustLoadOnce && !configLoaded {
		// Read from remote config the first time
		newCnf, err := LoadConfig(etcdClient)
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
				newCnf, err := LoadConfig(etcdClient)
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
func LoadConfig(cli *clientv3.Client) (*Config, error) {
	// Read from remote config the first time
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	resp, err := cli.Get(ctx, etcdConfigPath)
	cancel()
	if err != nil {
		return nil, err
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
        	// ETCD config
        	etcdClientConfig := clientv3.Config{
                	Endpoints:              strings.Split(theEndpoints, ","),
                	DialTimeout:            5 * time.Second,
                	TLS:                    tlsConfig,
        	}	
        	// ETCD client
        	etcdClient, err := clientv3.New(etcdClientConfig)
        	if err != nil {
                	return nil, err
        	}	
		return etcdClient, nil
	} else {
        	// ETCD config
        	etcdClientConfig := clientv3.Config{
                	Endpoints:              strings.Split(theEndpoints, ","),
                	DialTimeout:            5 * time.Second,
        	}	
        	// ETCD client
        	etcdClient, err := clientv3.New(etcdClientConfig)
        	if err != nil {
                	return nil, err
        	}
		return etcdClient, nil
	}
}
