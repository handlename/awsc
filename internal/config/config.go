package config

import (
	goconfig "github.com/kayac/go-config"
)

type Config struct {
	// patterns is highlight patterns for aws profile.
	// If profile name matches any of patterns, awsc outputs and hilights that
	Patterns []string `yaml:"patterns,omitempty"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		return &Config{}, nil
	}

	c := &Config{}
	if err := goconfig.LoadWithEnv(c, path); err != nil {
		return nil, err
	}

	return c, nil
}
