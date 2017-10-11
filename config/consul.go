package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/RichardKnop/go-oauth2-server/log"
	"github.com/hashicorp/consul/api"
)

var (
	consulEndpoint                              = "http://localhost:8500"
	consulCertFile, consulKeyFile, consulCaFile string
	consulConfigPath                            = "/config/go_oauth2_server.json"
)

type consulBackend struct{}

func (b *consulBackend) InitConfigBackend() {
	// Overwrite default values with environment variables if they are set
	if os.Getenv("CONSUL_ENDPOINT") != "" {
		consulEndpoint = os.Getenv("CONSUL_ENDPOINT")
	}
	if os.Getenv("CONSUL_CERT_FILE") != "" {
		consulCertFile = os.Getenv("CONSUL_CERT_FILE")
	}
	if os.Getenv("CONSUL_KEY_FILE") != "" {
		consulKeyFile = os.Getenv("CONSUL_KEY_FILE")
	}
	if os.Getenv("CONSUL_CA_FILE") != "" {
		consulCaFile = os.Getenv("CONSUL_CA_FILE")
	}
	if os.Getenv("CONSUL_CONFIG_PATH") != "" {
		consulConfigPath = os.Getenv("CONSUL_CONFIG_PATH")
	}
}

//LoadConfig gets the JSON from Consul and unmarshals it to the config object
func (b *consulBackend) LoadConfig() (*Config, error) {

	cli, err := newConsulClient(consulEndpoint, consulCertFile, consulKeyFile, consulCaFile)
	if err != nil {
		return nil, err
	}

	// Read from remote config the first time

	resp, _, err := cli.KV().Get(consulConfigPath, nil)

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("key not found: %s", consulConfigPath)
	}

	// Unmarshal the config JSON into the cnf object
	newCnf := new(Config)

	if err := json.Unmarshal(resp.Value, newCnf); err != nil {
		return nil, err
	}

	return newCnf, nil
}

// RefreshConfig sets config through the pointer so config actually gets refreshed
func (b *consulBackend) RefreshConfig(newCnf *Config) {
	*Cnf = *newCnf
}

func newConsulClient(theEndpoint, certFile, keyFile, caFile string) (*api.Client, error) {
	// Log the consul endpoint for debugging purposes
	log.INFO.Printf("CONSUL Endpoint: %s", theEndpoint)

	consulConfig := api.DefaultConfig()

	consulConfig.Address = theEndpoint

	// Optionally, configure TLS transport
	if certFile != "" && keyFile != "" && caFile != "" {

		tlsConfig, err := api.SetupTLSConfig(&api.TLSConfig{
			CertFile:           certFile,
			KeyFile:            keyFile,
			CAFile:             caFile,
			InsecureSkipVerify: true,
		})

		if err != nil {
			return nil, err
		}

		consulConfig.HttpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}

	}

	return api.NewClient(consulConfig)

}
