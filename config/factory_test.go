package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNewConfig(t *testing.T) {
	viper.Set("DATABASE_TYPE", "a")
	viper.Set("DATABASE_HOST", "b")
	viper.Set("DATABASE_PORT", 1)
	viper.Set("DATABASE_USER", "c")
	viper.Set("DATABASE_PASSWORD", "d")
	viper.Set("DATABASE_NAME", "e")
	viper.Set("ACCESS_TOKEN_LIFETIME", 2)
	viper.Set("REFRESH_TOKEN_LIFETIME", 3)

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
