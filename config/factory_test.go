package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	viper.Set("database_type", "a")
	viper.Set("database_host", "b")
	viper.Set("database_port", 1)
	viper.Set("database_user", "c")
	viper.Set("database_password", "d")
	viper.Set("database_name", "e")
	viper.Set("access_token_lifetime", 2)
	viper.Set("refresh_token_lifetime", 3)

	cnf := NewConfig()

	assert.Equal(t, "a", cnf.Database.Type, "cnf.Database.Type should be a")

	assert.Equal(t, "b", cnf.Database.Host, "cnf.Database.Host should be b")

	assert.Equal(t, 1, cnf.Database.Port, "cnf.Database.Port should be 1")

	assert.Equal(t, "c", cnf.Database.User, "cnf.Database.User should be c")

	assert.Equal(t, "d", cnf.Database.Password, "cnf.Database.Password should be d")

	assert.Equal(t, "e", cnf.Database.DatabaseName, "cnf.Database.DatabaseName should be e")

	assert.Equal(t, 2, cnf.AccessTokenLifetime, "cnf.AccessTokenLifetime should be 2")

	assert.Equal(t, 3, cnf.RefreshTokenLifetime, "cnf.RefreshTokenLifetime should be 3")
}
