package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/go-oauth2-server/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/pkg/transport"
	"golang.org/x/net/context"
)

var (
	etcdEndpoints                         = "http://localhost:2379"
	etcdCertFile, etcdKeyFile, etcdCaFile string
	etcdConfigPath                        = "/config/go_oauth2_server.json"
)

type etcdBackend struct{}

func (b *etcdBackend) InitConfigBackend() {
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

// LoadConfig gets the JSON from ETCD and unmarshals it to the config object
func (b *etcdBackend) LoadConfig() (*Config, error) {

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
func (b *etcdBackend) RefreshConfig(newCnf *Config) {
	*Cnf = *newCnf
}

func newEtcdClient(theEndpoints, certFile, keyFile, caFile string) (*clientv3.Client, error) {
	// Log the etcd endpoint for debugging purposes
	log.INFO.Printf("ETCD Endpoints: %s", theEndpoints)

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
