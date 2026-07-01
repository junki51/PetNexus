package database

import (
	"testing"

	"github.com/phonlakitz/petnexus-backend/internal/config"
)

func TestBuildPostgresDSNPrefersDatabaseURL(t *testing.T) {
	const databaseURL = "postgresql://render-user:secret@render-host/petnexus?sslmode=require"
	cfg := config.Config{
		DatabaseURL: databaseURL,
		DBHost:      "localhost",
		DBPort:      "5432",
		DBUser:      "postgres",
		DBPassword:  "postgres",
		DBName:      "petnexus",
		DBSSLMode:   "disable",
	}

	if got := buildPostgresDSN(cfg); got != databaseURL {
		t.Fatalf("buildPostgresDSN() = %q, want DATABASE_URL", got)
	}
}

func TestBuildPostgresDSNFallsBackToLocalConfig(t *testing.T) {
	cfg := config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "petnexus",
		DBSSLMode:  "disable",
	}
	want := "host=localhost port=5432 user=postgres password=postgres dbname=petnexus sslmode=disable"

	if got := buildPostgresDSN(cfg); got != want {
		t.Fatalf("buildPostgresDSN() = %q, want %q", got, want)
	}
}

func TestBuildPostgresDSNIgnoresWhitespaceDatabaseURL(t *testing.T) {
	cfg := config.Config{
		DatabaseURL: "   ",
		DBHost:      "localhost",
		DBPort:      "5432",
		DBUser:      "postgres",
		DBPassword:  "postgres",
		DBName:      "petnexus",
		DBSSLMode:   "disable",
	}
	want := "host=localhost port=5432 user=postgres password=postgres dbname=petnexus sslmode=disable"

	if got := buildPostgresDSN(cfg); got != want {
		t.Fatalf("buildPostgresDSN() = %q, want local DSN %q", got, want)
	}
}
