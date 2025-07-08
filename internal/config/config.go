package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration.
type Config struct {
	DefaultTags    []string               `json:"defaultTags"`
	OutputFormats  map[string]string      `json:"outputFormats"`
	FragmentsDir   string                 `json:"fragmentsDir,omitempty"`
	CustomSettings map[string]interface{} `json:"customSettings,omitempty"`
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	return &Config{
		DefaultTags: []string{},
		OutputFormats: map[string]string{
			"opencode": "AGENTS.md",
			"gemini":   "GEMINI.md",
		},
		FragmentsDir:   "",
		CustomSettings: make(map[string]interface{}),
	}
}

// GetConfigDir returns the configuration directory path.
func GetConfigDir() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}

		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, ".ctx"), nil
}

// GetFragmentsDir returns the fragments directory path.
func GetFragmentsDir(config *Config) (string, error) {
	if config.FragmentsDir != "" {
		return config.FragmentsDir, nil
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "fragments"), nil
}

// LoadConfig loads configuration from the specified file path.
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return nil, err
		}

		configPath = filepath.Join(configDir, "config.json")
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified file path.
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return err
		}

		configPath = filepath.Join(configDir, "config.json")
	}

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
