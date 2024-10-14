package config

import (
	goconfig "github.com/kayac/go-config"
)

type Config struct {
	// Rules is slice of caution display rule.
	Rules []Rule `yaml:"rules,omitempty"`

	// Template is output template
	// If empty, default template is used.
	// Available variables are:
	// - .Profile: AWS profile name
	// - .Region: AWS region
	// - .ID: AWS account ID (if ExtraInfo is true)
	// - .UserID: AWS user ID (if ExtraInfo is true)
	// - .Arn: AWS ARN (if ExtraInfo is true)
	// - .Now: Timestamp at command executed, formatted by `TimeFormat`
	Template string

	// ExtraInfo is flag to show extra information.
	// If true, awsc makes one or more API calls to AWS to gather extra data.
	ExtraInfo bool `yaml:"extra_info,omitempty"`

	// TimeFormat is format for time in output
	// Default is "2006-01-02 15:04:05"
	TimeFormat string `yaml:"time_format,omitempty"`
}

type Rule struct {
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
