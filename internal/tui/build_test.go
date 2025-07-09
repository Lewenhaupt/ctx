package tui

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Lewenhaupt/ctx/internal/config"
	"github.com/Lewenhaupt/ctx/internal/parser"
)

func TestDetermineSelectedTags(t *testing.T) {
	tests := []struct {
		name         string
		opts         *BuildOptions
		cfg          *config.Config
		fragments    []parser.Fragment
		expectedTags []string
		expectError  bool
	}{
		{
			name: "use provided tags",
			opts: &BuildOptions{
				Tags: []string{"typescript", "rust"},
			},
			cfg: &config.Config{
				DefaultTags: []string{"go", "python"},
			},
			fragments: []parser.Fragment{
				{Tags: []string{"typescript", "go"}},
				{Tags: []string{"rust", "python"}},
			},
			expectedTags: []string{"typescript", "rust"},
			expectError:  false,
		},
		{
			name: "use default tags in non-interactive mode",
			opts: &BuildOptions{
				NonInteractive: true,
			},
			cfg: &config.Config{
				DefaultTags: []string{"go", "python"},
			},
			fragments: []parser.Fragment{
				{Tags: []string{"typescript", "go"}},
				{Tags: []string{"rust", "python"}},
			},
			expectedTags: []string{"go", "python"},
			expectError:  false,
		},
		{
			name: "no tags found in fragments",
			opts: &BuildOptions{},
			cfg: &config.Config{
				DefaultTags: []string{},
			},
			fragments:   []parser.Fragment{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := determineSelectedTags(tt.opts, tt.cfg, tt.fragments)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expectedTags) {
				t.Errorf("Expected tags %v, got %v", tt.expectedTags, result)
			}
		})
	}
}

func TestLoadConfigAndFragments(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Set XDG_CONFIG_HOME to our temp directory
	oldXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		if oldXDGConfigHome != "" {
			_ = os.Setenv("XDG_CONFIG_HOME", oldXDGConfigHome)
		} else {
			_ = os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	_ = os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".ctx")

	err := os.MkdirAll(configDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := `{
		"defaultTags": ["typescript", "rust"],
		"outputFormats": {
			"test": "TEST.md"
		}
	}`
	configPath := filepath.Join(configDir, "config.json")

	err = os.WriteFile(configPath, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create fragments directory and files
	fragmentsDir := filepath.Join(configDir, "fragments")

	err = os.MkdirAll(fragmentsDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create fragments directory: %v", err)
	}

	fragmentContent := `---
ctx-tags: typescript, testing
---

# Test Fragment
This is a test fragment.
`
	fragmentPath := filepath.Join(fragmentsDir, "test.md")

	err = os.WriteFile(fragmentPath, []byte(fragmentContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create fragment file: %v", err)
	}

	// Test loading config and fragments
	cfg, fragments, err := loadConfigAndFragments("", false)
	if err != nil {
		t.Fatalf("loadConfigAndFragments failed: %v", err)
	}

	// Verify config
	expectedDefaultTags := []string{"typescript", "rust"}
	if !reflect.DeepEqual(cfg.DefaultTags, expectedDefaultTags) {
		t.Errorf("Expected default tags %v, got %v", expectedDefaultTags, cfg.DefaultTags)
	}

	// Verify fragments
	if len(fragments) != 1 {
		t.Errorf("Expected 1 fragment, got %d", len(fragments))
	}

	if len(fragments) > 0 {
		expectedTags := []string{"typescript", "testing"}
		if !reflect.DeepEqual(fragments[0].Tags, expectedTags) {
			t.Errorf("Expected fragment tags %v, got %v", expectedTags, fragments[0].Tags)
		}
	}
}

func TestHandleFileOverwrite(t *testing.T) {
	tests := []struct {
		name           string
		opts           *BuildOptions
		filename       string
		format         string
		expectedAction string
	}{
		{
			name: "non-interactive mode always overwrites",
			opts: &BuildOptions{
				NonInteractive: true,
			},
			filename:       "test.md",
			format:         "opencode",
			expectedAction: "overwrite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := handleFileOverwrite(tt.opts, tt.filename, tt.format)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if action != tt.expectedAction {
				t.Errorf("Expected action %s, got %s", tt.expectedAction, action)
			}
		})
	}
}

func TestDetermineOutputFormats(t *testing.T) {
	tests := []struct {
		name            string
		opts            *BuildOptions
		cfg             *config.Config
		expectedFormats []string
		expectedFiles   []string
		expectError     bool
	}{
		{
			name: "stdout option",
			opts: &BuildOptions{
				Stdout: true,
			},
			cfg:             &config.Config{},
			expectedFormats: []string{"stdout"},
			expectedFiles:   nil,
			expectError:     false,
		},
		{
			name: "provided output formats",
			opts: &BuildOptions{
				OutputFormats: []string{"opencode", "gemini"},
			},
			cfg:             &config.Config{},
			expectedFormats: []string{"opencode", "gemini"},
			expectedFiles:   nil,
			expectError:     false,
		},
		{
			name: "custom output file",
			opts: &BuildOptions{
				OutputFile: "custom.md",
			},
			cfg:             &config.Config{},
			expectedFormats: []string{"custom"},
			expectedFiles:   []string{"custom.md"},
			expectError:     false,
		},
		{
			name: "non-interactive with configured formats",
			opts: &BuildOptions{
				NonInteractive: true,
			},
			cfg: &config.Config{
				OutputFormats: map[string]string{
					"opencode": "AGENTS.md",
					"gemini":   "GEMINI.md",
				},
			},
			expectedFormats: []string{"opencode", "gemini"},
			expectedFiles:   nil,
			expectError:     false,
		},
		{
			name: "non-interactive with no configured formats",
			opts: &BuildOptions{
				NonInteractive: true,
			},
			cfg: &config.Config{
				OutputFormats: map[string]string{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formats, files, err := determineOutputFormats(tt.opts, tt.cfg)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// For non-interactive mode with multiple formats, the order might vary
			// so we need to check if all expected formats are present
			if tt.name == "non-interactive with configured formats" {
				if len(formats) != len(tt.expectedFormats) {
					t.Errorf("Expected %d formats, got %d", len(tt.expectedFormats), len(formats))
				}

				formatMap := make(map[string]bool)
				for _, format := range formats {
					formatMap[format] = true
				}

				for _, expected := range tt.expectedFormats {
					if !formatMap[expected] {
						t.Errorf("Expected format %s not found in result", expected)
					}
				}
			} else if !reflect.DeepEqual(formats, tt.expectedFormats) {
				t.Errorf("Expected formats %v, got %v", tt.expectedFormats, formats)
			}

			if !reflect.DeepEqual(files, tt.expectedFiles) {
				t.Errorf("Expected files %v, got %v", tt.expectedFiles, files)
			}
		})
	}
}
