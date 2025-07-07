package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationBuild(t *testing.T) {
	// Create temporary directory for test fragments
	tmpDir := t.TempDir()
	fragmentsDir := filepath.Join(tmpDir, "fragments")
	err := os.MkdirAll(fragmentsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create fragments directory: %v", err)
	}

	// Create test fragments
	fragments := map[string]string{
		"typescript.md": `---
ctx-tags: typescript, frontend
---

# TypeScript Guidelines
Use strict mode.`,
		"rust.md": `---
ctx-tags: rust, systems
---

# Rust Guidelines
Use the borrow checker.`,
		"general.md": `---
ctx-tags: general, coding
---

# General Guidelines
Write clean code.`,
	}

	for filename, content := range fragments {
		err := os.WriteFile(filepath.Join(fragmentsDir, filename), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create fragment %s: %v", filename, err)
		}
	}

	// Create config file in .ctx directory
	ctxDir := filepath.Join(tmpDir, ".ctx")
	err = os.MkdirAll(ctxDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .ctx directory: %v", err)
	}

	configContent := `{
		"default_tags": ["general"],
		"fragments_dir": "` + fragmentsDir + `"
	}`
	configPath := filepath.Join(ctxDir, "config.json")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Set XDG_CONFIG_HOME to our test directory
	oldXDGConfig := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if oldXDGConfig != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDGConfig)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:           "build with typescript tag",
			args:           []string{"build", "--tags", "typescript", "--non-interactive"},
			expectedOutput: []string{"# TypeScript Guidelines", "Use strict mode."},
			expectError:    false,
		},
		{
			name:           "build with multiple tags",
			args:           []string{"build", "--tags", "typescript,rust", "--non-interactive"},
			expectedOutput: []string{"# TypeScript Guidelines", "# Rust Guidelines"},
			expectError:    false,
		},
		{
			name:        "build with nonexistent tag",
			args:        []string{"build", "--tags", "nonexistent", "--non-interactive"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the binary
			cmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "ctx"), "./cmd/ctx")
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Failed to build binary: %v", err)
			}

			// Run the command
			cmd = exec.Command(filepath.Join(tmpDir, "ctx"), tt.args...)
			cmd.Env = append(os.Environ(), "XDG_CONFIG_HOME="+tmpDir)
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
			}

			outputStr := string(output)
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, outputStr)
				}
			}
		})
	}
}