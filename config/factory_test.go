package config

import (
	"testing"

	"github.com/spf13/viper"
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

	if cnf.Database.Type != "a" {
		t.Errorf("cnf.Database.Type = %s, want a", cnf.Database.Type)
	}

	if cnf.Database.Host != "b" {
		t.Errorf("cnf.Database.Host = %s, want b", cnf.Database.Host)
	}

	if cnf.Database.Port != 1 {
		t.Errorf("cnf.Database.Port = %d, want 1", cnf.Database.Port)
	}

	if cnf.Database.User != "c" {
		t.Errorf("cnf.Database.User = %s, want c", cnf.Database.User)
	}

	if cnf.Database.Password != "d" {
		t.Errorf("cnf.Database.Password = %s, want d", cnf.Database.Password)
	}

	if cnf.Database.DatabaseName != "e" {
		t.Errorf("cnf.Database.DatabaseName = %s, want e", cnf.Database.Password)
	}

	if cnf.AccessTokenLifetime != 2 {
		t.Errorf("cnf.AccessTokenLifetime = %d, want 2", cnf.AccessTokenLifetime)
	}

	if cnf.RefreshTokenLifetime != 3 {
		t.Errorf("cnf.RefreshTokenLifetime = %d, want 3", cnf.RefreshTokenLifetime)
	}
}
