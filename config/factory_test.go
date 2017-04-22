package config_test

import (
	"testing"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	cnf *config.Config
}

func (b *Foo) LoadConfig() (*config.Config, error) {
	return nil, nil
}

func (b *Foo) RefreshConfig(newCnf *config.Config) {
	b.cnf = newCnf
}

func (b *Foo) InitConfigBackend() {}

func TestConfigReloading(t *testing.T) {

	config.Cnf.Oauth.AuthCodeLifetime = 123
	foo := &Foo{cnf: config.Cnf}
	assert.Equal(t, 123, foo.cnf.Oauth.AuthCodeLifetime)
	newCnf := &config.Config{Oauth: config.OauthConfig{AuthCodeLifetime: 9999}}
	assert.Equal(t, 123, foo.cnf.Oauth.AuthCodeLifetime)
	foo.RefreshConfig(newCnf)
	assert.Equal(t, 9999, foo.cnf.Oauth.AuthCodeLifetime)
}
