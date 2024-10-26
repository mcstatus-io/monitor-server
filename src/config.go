package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// DefaultConfig is the default configuration values used by the application.
	DefaultConfig *Config = &Config{
		Environment:    "production",
		CycleInterval:  time.Minute,
		PingServerHost: "http://localhost:3001",
		AuthToken:      "",
		MongoDB:        "mongodb://127.0.0.1:27017/mcstatus",
	}
)

// Config represents the application configuration.
type Config struct {
	Environment    string        `yaml:"environment"`
	CycleInterval  time.Duration `yaml:"cycle_interval"`
	PingServerHost string        `yaml:"ping_server_host"`
	AuthToken      string        `yaml:"auth_token"`
	MongoDB        string        `yaml:"mongodb"`
}

// ReadFile reads the configuration from the given file and overrides values using environment variables.
func (c *Config) ReadFile(file string) error {
	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	return c.overrideWithEnvVars()
}

// WriteFile writes the configuration values to a file.
func (c *Config) WriteFile(file string) error {
	data, err := yaml.Marshal(c)

	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0777)
}

func (c *Config) overrideWithEnvVars() error {
	if value := os.Getenv("ENVIRONMENT"); value != "" {
		c.Environment = value
	}

	if value := os.Getenv("MONGO_URL"); value != "" {
		c.MongoDB = value
	}

	return nil
}
