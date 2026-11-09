package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	for name, value := range map[string]string{
		"DB_DRIVER":            "mysql",
		"DB_HOST":              "127.0.0.1",
		"DB_PORT":              "3306",
		"DB_USER":              "test-user",
		"DB_NAME":              "orders",
		"WEB_SERVER_PORT":      ":8000",
		"GRPC_SERVER_PORT":     "50051",
		"GRAPHQL_SERVER_PORT": "8080",
	} {
		t.Setenv(name, value)
	}
}

func TestLoadConfigWithoutDotenvUsesEnvironment(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected config error without .env: %v", err)
	}
	if cfg.DBHost != "127.0.0.1" || cfg.DBDriver != "mysql" || cfg.WebServerPort != ":8000" {
		t.Fatalf("environment settings not loaded correctly: %+v", cfg)
	}
}

func TestLoadConfigEnvironmentOverridesDotenv(t *testing.T) {
	dir := t.TempDir()
	body := "DB_HOST=from-file\\nDB_DRIVER=mysql\\nWEB_SERVER_PORT=:8000\\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	setRequiredEnv(t)
	t.Setenv("DB_HOST", "from-environment")

	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBHost != "from-environment" || cfg.DBDriver != "mysql" || cfg.WebServerPort != ":8000" {
		t.Fatalf("unexpected configuration precedence: %+v", cfg)
	}
}

func TestLoadConfigReturnsErrorForInvalidDotenvFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".env"), 0700); err != nil {
		t.Fatal(err)
	}

	setRequiredEnv(t)
	cfg, err := LoadConfig(dir)
	if err == nil || cfg != nil {
		t.Fatalf("unreadable .env must return error rather than panic: cfg=%+v err=%v", cfg, err)
	}
}

func TestLoadConfigReportsMissingRequiredKey(t *testing.T) {
	// Do not call setRequiredEnv here: an empty temporary directory must
	// be rejected rather than returning a partly initialized config.
	t.Setenv("DB_DRIVER", "")
	cfg, err := LoadConfig(t.TempDir())
	if err == nil || cfg != nil {
		t.Fatalf("expected missing configuration error: cfg=%+v err=%v", cfg, err)
	}
}
