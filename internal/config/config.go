package config

import (
	goconfig "github.com/kayac/go-config"
)

type Config struct {
	// Patterns is highlight patterns for aws profile.
	// If profile name matches any of patterns, awsc outputs and hilights that
	Patterns []Pattern `yaml:"patterns,omitempty"`

	// AdditionalInfo is flag to show additional information.
	// If true, awsc shows additional information, but it calls additional AWS API.
	AdditionalInfo bool `yaml:"additional_info,omitempty"`

	// TimeFormat is format for time in output
	// Default is "2006-01-02 15:04:05"
	TimeFormat string `yaml:"time_format,omitempty"`
}

type Pattern struct {
	// Expression is regexp for AWS profile name to match
	Expression string `yaml:"expression,omitempty"`

	// Color is highlight color for matched profile name
	Color string `yaml:"color,omitempty"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		return &Config{}, nil
	}

	c := &Config{}
	if err := goconfig.LoadWithEnv(c, path); err != nil {
		return nil, err
	}

	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}

	return c, nil
}
