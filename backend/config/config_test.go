package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Unsetenv("DB_PATH")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("PORT")

	cfg := LoadConfig()

	if cfg.DBPath != "./data/go4movies.db" {
		t.Errorf("expected default DBPath, got %q", cfg.DBPath)
	}
	if cfg.JWTSecret != "change-me-in-production" {
		t.Errorf("expected default JWTSecret, got %q", cfg.JWTSecret)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default Port 8080, got %q", cfg.Port)
	}
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	t.Setenv("DB_PATH", "/tmp/test.db")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("PORT", "9090")

	cfg := LoadConfig()

	if cfg.DBPath != "/tmp/test.db" {
		t.Errorf("expected /tmp/test.db, got %q", cfg.DBPath)
	}
	if cfg.JWTSecret != "super-secret" {
		t.Errorf("expected super-secret, got %q", cfg.JWTSecret)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected 9090, got %q", cfg.Port)
	}
}

func TestGenENV_Fallback(t *testing.T) {
	os.Unsetenv("NONEXISTENT_KEY")
	val := genENV("NONEXISTENT_KEY", "default-value")
	if val != "default-value" {
		t.Errorf("expected fallback, got %q", val)
	}
}

func TestGenENV_EnvSet(t *testing.T) {
	t.Setenv("TEST_KEY", "from-env")
	val := genENV("TEST_KEY", "default")
	if val != "from-env" {
		t.Errorf("expected from-env, got %q", val)
	}
}
