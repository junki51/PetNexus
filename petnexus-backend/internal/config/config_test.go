package config

import "testing"

func TestLoadReadsDatabaseURL(t *testing.T) {
	const databaseURL = "postgresql://render-user:secret@render-host/petnexus?sslmode=require"
	t.Setenv("DATABASE_URL", databaseURL)

	cfg := Load()

	if cfg.DatabaseURL != databaseURL {
		t.Fatalf("Load().DatabaseURL = %q, want %q", cfg.DatabaseURL, databaseURL)
	}
}

func TestLoadDefaultsDatabaseURLToEmpty(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	cfg := Load()

	if cfg.DatabaseURL != "" {
		t.Fatalf("Load().DatabaseURL = %q, want empty string", cfg.DatabaseURL)
	}
}
