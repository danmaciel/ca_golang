package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigWithoutDotenvUsesEnvironment(t *testing.T) {
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_DRIVER", "mysql")
	t.Setenv("WEB_SERVER_PORT", ":8000")

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

	cfg, err := LoadConfig(dir)
	if err == nil || cfg != nil {
		t.Fatalf("unreadable .env must return error rather than panic: cfg=%+v err=%v", cfg, err)
	}
}
