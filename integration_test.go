package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationBuild(t *testing.T) {
	tmpDir := setupTestEnvironment(t)
	tests := getTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runIntegrationTest(t, tmpDir, tt)
		})
	}
}

func setupTestEnvironment(t *testing.T) string {
	tmpDir := t.TempDir()

	createTestFragments(t, tmpDir)
	createTestConfig(t, tmpDir)
	setupEnvironmentVariables(t, tmpDir)

	return tmpDir
}

func createTestFragments(t *testing.T, tmpDir string) {
	fragmentsDir := filepath.Join(tmpDir, "fragments")

	err := os.MkdirAll(fragmentsDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create fragments directory: %v", err)
	}

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
		err := os.WriteFile(filepath.Join(fragmentsDir, filename), []byte(content), 0o600)
		if err != nil {
			t.Fatalf("Failed to create fragment %s: %v", filename, err)
		}
	}
}

func createTestConfig(t *testing.T, tmpDir string) {
	fragmentsDir := filepath.Join(tmpDir, "fragments")
	ctxDir := filepath.Join(tmpDir, ".ctx")

	err := os.MkdirAll(ctxDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create .ctx directory: %v", err)
	}

	configContent := `{
		"default_tags": ["general"],
		"fragments_dir": "` + fragmentsDir + `",
		"output_formats": {
			"opencode": "AGENTS.md"
		}
	}`
	configPath := filepath.Join(ctxDir, "config.json")

	err = os.WriteFile(configPath, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
}

func setupEnvironmentVariables(t *testing.T, tmpDir string) {
	oldXDGConfig := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		if oldXDGConfig != "" {
			if err := os.Setenv("XDG_CONFIG_HOME", oldXDGConfig); err != nil {
				t.Errorf("Failed to restore XDG_CONFIG_HOME: %v", err)
			}
		} else {
			if err := os.Unsetenv("XDG_CONFIG_HOME"); err != nil {
				t.Errorf("Failed to unset XDG_CONFIG_HOME: %v", err)
			}
		}
	}()

	if err := os.Setenv("XDG_CONFIG_HOME", tmpDir); err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}
}

func getTestCases() []struct {
	name           string
	args           []string
	expectedOutput []string
	expectError    bool
} {
	return []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:           "build with typescript tag",
			args:           []string{"build", "--tags", "typescript", "--non-interactive", "--stdout"},
			expectedOutput: []string{"# TypeScript Guidelines", "Use strict mode."},
			expectError:    false,
		},
		{
			name:           "build with multiple tags",
			args:           []string{"build", "--tags", "typescript,rust", "--non-interactive", "--stdout"},
			expectedOutput: []string{"# TypeScript Guidelines", "# Rust Guidelines"},
			expectError:    false,
		},
		{
			name:        "build with nonexistent tag",
			args:        []string{"build", "--tags", "nonexistent", "--non-interactive", "--stdout"},
			expectError: true,
		},
	}
}

func runIntegrationTest(t *testing.T, tmpDir string, tt struct {
	name           string
	args           []string
	expectedOutput []string
	expectError    bool
},
) {
	buildBinary(t, tmpDir)
	output, err := runCommand(t, tmpDir, tt.args)

	validateTestResult(t, output, err, tt.expectedOutput, tt.expectError)
}

func buildBinary(t *testing.T, tmpDir string) {
	cmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "ctx"), "./cmd/ctx")

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
}

func runCommand(t *testing.T, tmpDir string, args []string) ([]byte, error) {
	cmd := exec.Command(filepath.Join(tmpDir, "ctx"), args...)
	cmd.Env = append(os.Environ(), "XDG_CONFIG_HOME="+tmpDir)

	return cmd.CombinedOutput()
}

func validateTestResult(t *testing.T, output []byte, err error, expectedOutput []string, expectError bool) {
	if expectError {
		if err == nil {
			t.Error("Expected error, got nil")
		}

		return
	}

	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	for _, expected := range expectedOutput {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Expected output to contain %q, got: %s", expected, outputStr)
		}
	}
}
