package configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

// LoadConfig reads an optional local .env and lets process environment values
// override it. Each call owns its Viper instance so tests do not share state.
func LoadConfig(path string) (*conf, error) {
	v := viper.New()
	v.SetConfigFile(filepath.Join(path, ".env"))
	v.SetConfigType("env")
	v.AutomaticEnv()

	keys := []string{
		"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD",
		"DB_NAME", "WEB_SERVER_PORT", "GRPC_SERVER_PORT", "GRAPHQL_SERVER_PORT",
	}
	for _, key := range keys {
		if err := v.BindEnv(key); err != nil {
			return nil, fmt.Errorf("bind configuration %s: %w", key, err)
		}
	}

	if err := v.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read local configuration: %w", err)
	}

	// Fail early with variable names only, never print configured credentials.
	for _, required := range []string{
		"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_USER", "DB_NAME",
		"WEB_SERVER_PORT", "GRPC_SERVER_PORT", "GRAPHQL_SERVER_PORT",
	} {
		if strings.TrimSpace(v.GetString(required)) == "" {
			return nil, fmt.Errorf("missing required configuration: %s", required)
		}
	}

	var cfg conf
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decode configuration: %w", err)
	}
	return &cfg, nil
}
