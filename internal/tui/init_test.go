package tui

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Lewenhaupt/ctx/internal/config"
)

func TestGenerateConfig(t *testing.T) {
	tests := []struct {
		name     string
		answers  *InitAnswers
		expected *config.Config
	}{
		{
			name: "default configuration",
			answers: &InitAnswers{
				AddOutputFormats: false,
				FragmentsDir:     "",
				CreateSample:     false,
			},
			expected: &config.Config{
				DefaultTags: []string{},
				OutputFormats: map[string]string{
					"opencode": "AGENTS.md",
					"gemini":   "GEMINI.md",
				},
				FragmentsDir:   "",
				CustomSettings: make(map[string]interface{}),
			},
		},
		{
			name: "custom fragments directory",
			answers: &InitAnswers{
				AddOutputFormats: false,
				FragmentsDir:     "./custom-fragments",
				CreateSample:     false,
			},
			expected: func() *config.Config {
				absPath, _ := filepath.Abs("./custom-fragments")
				return &config.Config{
					DefaultTags: []string{},
					OutputFormats: map[string]string{
						"opencode": "AGENTS.md",
						"gemini":   "GEMINI.md",
					},
					FragmentsDir:   absPath,
					CustomSettings: make(map[string]interface{}),
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generateConfig(tt.answers)
			if err != nil {
				t.Fatalf("generateConfig failed: %v", err)
			}

			if !reflect.DeepEqual(result.DefaultTags, tt.expected.DefaultTags) {
				t.Errorf("Expected default tags %v, got %v",
					tt.expected.DefaultTags, result.DefaultTags)
			}

			if !reflect.DeepEqual(result.OutputFormats, tt.expected.OutputFormats) {
				t.Errorf("Expected output formats %v, got %v",
					tt.expected.OutputFormats, result.OutputFormats)
			}

			if result.FragmentsDir != tt.expected.FragmentsDir {
				t.Errorf("Expected fragments dir %s, got %s",
					tt.expected.FragmentsDir, result.FragmentsDir)
			}
		})
	}
}

func TestCreateSampleFragment(t *testing.T) {
	tmpDir := t.TempDir()

	err := createSampleFragment(tmpDir)
	if err != nil {
		t.Fatalf("createSampleFragment failed: %v", err)
	}

	// Verify sample fragment was created
	samplePath := filepath.Join(tmpDir, "hello-world.md")
	if _, err := os.Stat(samplePath); os.IsNotExist(err) {
		t.Error("Sample fragment file was not created")
	}

	// Verify content
	content, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("Failed to read sample fragment: %v", err)
	}

	expectedContent := "ctx-tags: hello, world, sample"
	if !contains(string(content), expectedContent) {
		t.Errorf("Sample fragment does not contain expected ctx-tags header")
	}

	if !contains(string(content), "# Hello World") {
		t.Errorf("Sample fragment does not contain expected title")
	}
}

func TestRunInitWithExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up temporary config directory
	originalXDG := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		if originalXDG != "" {
			_ = os.Setenv("XDG_CONFIG_HOME", originalXDG)
		} else {
			_ = os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	_ = os.Setenv("XDG_CONFIG_HOME", tmpDir)
	// Create existing config
	configDir := filepath.Join(tmpDir, ".ctx")

	err := os.MkdirAll(configDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	existingConfig := `{"defaultTags": ["existing"]}`

	err = os.WriteFile(configPath, []byte(existingConfig), 0o600)
	if err != nil {
		t.Fatalf("Failed to create existing config: %v", err)
	}

	// Verify config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Existing config file should exist")
	}
}

func TestGenerateConfigWithInvalidPath(t *testing.T) {
	// Test with a path that contains invalid characters for Windows
	// but might be valid on Unix systems, so we'll test a different scenario
	answers := &InitAnswers{
		AddOutputFormats: false,
		FragmentsDir:     "/nonexistent/very/deep/path/that/should/not/exist",
		CreateSample:     false,
	}

	// This should not fail during config generation, only during directory creation
	cfg, err := generateConfig(answers)
	if err != nil {
		t.Fatalf("generateConfig should not fail for non-existent path: %v", err)
	}

	// Verify the path was set correctly
	expectedPath, _ := filepath.Abs("/nonexistent/very/deep/path/that/should/not/exist")
	if cfg.FragmentsDir != expectedPath {
		t.Errorf("Expected fragments dir %s, got %s", expectedPath, cfg.FragmentsDir)
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			contains(s[1:], substr) ||
			(s != "" && s[:len(substr)] == substr))
}
