package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if len(config.DefaultTags) != 0 {
		t.Errorf("Expected empty default tags, got %v", config.DefaultTags)
	}

	expectedFormats := map[string]string{
		"opencode": "AGENTS.md",
		"gemini":   "GEMINI.md",
	}

	if !reflect.DeepEqual(config.OutputFormats, expectedFormats) {
		t.Errorf("Expected output formats %v, got %v", expectedFormats, config.OutputFormats)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name: "valid config",
			configContent: `{
				"default_tags": ["typescript", "rust"],
				"output_formats": {
					"custom": "CUSTOM.md"
				}
			}`,
			expectedConfig: &Config{
				DefaultTags: []string{"typescript", "rust"},
				OutputFormats: map[string]string{
					"custom": "CUSTOM.md",
				},
				CustomSettings: make(map[string]interface{}),
			},
			expectError: false,
		},
		{
			name:          "invalid json",
			configContent: `{invalid json}`,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.json")

			if tt.configContent != "" {
				err := os.WriteFile(configPath, []byte(tt.configContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}
			}

			config, err := LoadConfig(configPath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectedConfig != nil {
				if !reflect.DeepEqual(config.DefaultTags, tt.expectedConfig.DefaultTags) {
					t.Errorf("Expected default tags %v, got %v",
						tt.expectedConfig.DefaultTags, config.DefaultTags)
				}

				if !reflect.DeepEqual(config.OutputFormats, tt.expectedConfig.OutputFormats) {
					t.Errorf("Expected output formats %v, got %v",
						tt.expectedConfig.OutputFormats, config.OutputFormats)
				}
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	config := &Config{
		DefaultTags: []string{"test"},
		OutputFormats: map[string]string{
			"test": "TEST.md",
		},
		CustomSettings: make(map[string]interface{}),
	}

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file was created and contains correct data
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var savedConfig Config
	err = json.Unmarshal(data, &savedConfig)
	if err != nil {
		t.Fatalf("Failed to parse saved config: %v", err)
	}

	if !reflect.DeepEqual(savedConfig.DefaultTags, config.DefaultTags) {
		t.Errorf("Expected default tags %v, got %v",
			config.DefaultTags, savedConfig.DefaultTags)
	}
}
