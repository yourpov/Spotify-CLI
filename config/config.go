package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"spotify-cli/internal/logx"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Load the config
func Load() (*Config, error) {
	paths := []string{
		filepath.Join("config", "config.json"),
		func() string {
			if base, err := os.UserConfigDir(); err == nil {
				return filepath.Join(base, "spotify-cli", "config.json")
			}
			return ""
		}(),
	}

	var lastErr error
	for _, p := range paths {
		if p == "" {
			continue
		}
		f, err := os.Open(p)
		if err != nil {
			lastErr = err
			continue
		}
		defer f.Close()

		var cfg Config
		if decErr := json.NewDecoder(f).Decode(&cfg); decErr != nil {
			return nil, logx.Err("decode %s: %w", p, decErr)
		}
		if cfg.ClientID == "" || cfg.ClientSecret == "" {
			return nil, logx.Err("%s missing client_id or client_secret", p)
		}
		return &cfg, nil
	}

	if lastErr == nil {
		lastErr = logx.Err("no config paths tried")
	}
	return nil, logx.Err("could not load config.json: %w", lastErr)
}
