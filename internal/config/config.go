package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server    ServerConfig    `toml:"server"`
	Auth      AuthConfig      `toml:"auth"`
	Storage   StorageConfig   `toml:"storage"`
	RateLimit RateLimitConfig `toml:"rate_limit"`
	Logging   LoggingConfig   `toml:"logging"`
}

type ServerConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	BaseURL string `toml:"base_url"`
	DataDir string `toml:"data_dir"`
	Dev     bool   `toml:"dev"`
}

type AuthConfig struct {
	JWTSecret        string `toml:"jwt_secret"`
	JWTExpiry        string `toml:"jwt_expiry"`
	BcryptCost       int    `toml:"bcrypt_cost"`
	RegistrationOpen bool   `toml:"registration_open"`
}

type StorageConfig struct {
	MaxObjectSize    string `toml:"max_object_size"`
	MaxSendSize      string `toml:"max_send_size"`
	MaxLoomsPerUser  int    `toml:"max_looms_per_user"`
}

type RateLimitConfig struct {
	Enabled      bool `toml:"enabled"`
	APIPerMinute int  `toml:"api_per_minute"`
}

type LoggingConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

func Defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Host:    "0.0.0.0",
			Port:    3000,
			BaseURL: "http://localhost:3000",
			DataDir: "./data",
		},
		Auth: AuthConfig{
			JWTExpiry:        "168h",
			BcryptCost:       12,
			RegistrationOpen: true,
		},
		Storage: StorageConfig{
			MaxObjectSize:   "100MB",
			MaxSendSize:     "500MB",
			MaxLoomsPerUser: 100,
		},
		RateLimit: RateLimitConfig{
			Enabled:      true,
			APIPerMinute: 300,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := Defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No config file — use defaults
			if cfg.Auth.JWTSecret == "" {
				cfg.Auth.JWTSecret = generateSecret()
			}
			return cfg, nil
		}
		return nil, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Auth.JWTSecret == "" {
		cfg.Auth.JWTSecret = generateSecret()
	}

	return cfg, nil
}

func generateSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
