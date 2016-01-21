package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

var (
	etcdHost     = "localhost"
	etcdPort     = "2379"
	etcdURL      string
	configPath   = "/config/go_oauth2_server.json"
	configLoaded bool
)

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

func init() {
	if os.Getenv("ETCD_HOST") != "" {
		etcdHost = os.Getenv("ETCD_HOST")
	}
	if os.Getenv("ETCD_PORT") != "" {
		etcdPort = os.Getenv("ETCD_PORT")
	}
	etcdURL = fmt.Sprintf("http://%s:%s", etcdHost, etcdPort)
}

// NewConfig loads configuration from etcd and returns *Config struct
// It also starts a goroutine in the background to keep config up-to-date
func NewConfig() *Config {
	if configLoaded {
		return cnf
	}

	// Construct the ETCD URL
	etcdHost := "localhost"
	if os.Getenv("ETCD_HOST") != "" {
		etcdHost = os.Getenv("ETCD_HOST")
	}
	etcdPort := "2379"
	if os.Getenv("ETCD_PORT") != "" {
		etcdPort = os.Getenv("ETCD_PORT")
	}
	etcdURL := fmt.Sprintf("http://%s:%s", etcdHost, etcdPort)
	log.Printf("ETCD URL: %s", etcdURL)

	// Initialise the SDK
	etcdClientConfig := client.Config{
		Endpoints: []string{etcdURL},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	etcdClient, err := client.New(etcdClientConfig)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(etcdClient)

	// Read from remote config the first time
	if err := loadConfig(kapi); err != nil {
		log.Fatal(err)
	}

	configLoaded = true
	log.Print("Successfully loaded config for the first time")

	// Open a goroutine to watch remote changes forever
	go func() {
		for {
			// Delay after each request
			time.Sleep(time.Second * 10)

			// Attempt to reload the config
			if err := loadConfig(kapi); err != nil {
				log.Print(err)
				return
			}

			log.Print("Successfully reloaded config")
		}
	}()

	return cnf
}

func loadConfig(kapi client.KeysAPI) error {
	// Read from remote config the first time
	resp, err := kapi.Get(context.Background(), configPath, nil)
	if err != nil {
		return err
	}

	// Unmarshal the config JSON into the cnf object
	if err := json.Unmarshal([]byte(resp.Node.Value), cnf); err != nil {
		return err
	}

	return nil
}
