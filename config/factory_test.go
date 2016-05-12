package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	cnf *Config
}

func TestConfigReloading(t *testing.T) {
	cnf.Oauth.AuthCodeLifetime = 123
	foo := &Foo{cnf: cnf}
	assert.Equal(t, 123, foo.cnf.Oauth.AuthCodeLifetime)
	newCnf := &Config{Oauth: OauthConfig{AuthCodeLifetime: 9999}}
	assert.Equal(t, 123, foo.cnf.Oauth.AuthCodeLifetime)
	refreshConfig(newCnf)
	assert.Equal(t, 9999, foo.cnf.Oauth.AuthCodeLifetime)
}
