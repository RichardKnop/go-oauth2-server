package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
	// Enable the remote features
	_ "github.com/spf13/viper/remote"
)

var cnf *Config

// NewConfig loads configuration from etcd and returns *Config struct
// It also starts a goroutine in the background to keep config up-to-date
func NewConfig() *Config {
	if cnf != nil {
		log.Print("AAAAAAA")
		log.Print(cnf)
		return cnf
	}

	runtimeViper := viper.New()
	runtimeViper.AddRemoteProvider(
		"etcd",
		"http://127.0.0.1:4001",
		"/config/go_oauth2_server.json",
	)
	// Because there is no file extension in a stream of bytes
	runtimeViper.SetConfigType("json")

	// Read from remote config the first time.
	if err := runtimeViper.ReadRemoteConfig(); err != nil {
		log.Fatal(err)
	}

	// Unmarshal config
	runtimeViper.Unmarshal(&cnf)

	// Open a goroutine to watch remote changes forever
	go func() {
		for {
			// Delay after each request
			time.Sleep(time.Second * 5)

			if err := runtimeViper.WatchRemoteConfig(); err != nil {
				log.Printf("Unable to read remote config: %v", err)
				continue
			}

			// Unmarshal config
			runtimeViper.Unmarshal(&cnf)
		}
	}()

	log.Print("BBBBBBB")
	log.Print(cnf)
	return cnf
}
