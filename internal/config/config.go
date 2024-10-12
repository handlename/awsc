package config

import (
	goconfig "github.com/kayac/go-config"
)

type Config struct {
	// Patterns is highlight patterns for aws profile.
	// If profile name matches any of patterns, awsc outputs and hilights that
	Patterns []Pattern `yaml:"patterns,omitempty"`

	// Template is output template
	// If empty, default template is used.
	// Available variables are:
	// - .Profile: AWS profile name
	// - .Region: AWS region
	// - .ID: AWS account ID (if AdditionalInfo is true)
	// - .UserID: AWS user ID (if AdditionalInfo is true)
	// - .Arn: AWS ARN (if AdditionalInfo is true)
	// - .Now: Timestamp at command executed, formatted by `TimeFormat`
	Template string

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

	// Change delimiters to avoid conflict with text/template's ones in `Template`
	loader := goconfig.New()
	loader.Delims("<<", ">>")

	c := &Config{}
	if err := loader.LoadWithEnv(c, path); err != nil {
		return nil, err
	}

	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}

	return c, nil
}
