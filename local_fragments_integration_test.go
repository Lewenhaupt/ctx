package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type integrationTestSetup struct {
	tmpDir           string
	ctxBinary        string
	originalWd       string
	oldXDGConfigHome string
}

func setupIntegrationTest(t *testing.T) *integrationTestSetup {
	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create temporary directory for integration test
	tmpDir := t.TempDir()

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Set XDG_CONFIG_HOME to a separate directory for global config
	globalConfigDir := filepath.Join(tmpDir, "global-config")
	oldXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	_ = os.Setenv("XDG_CONFIG_HOME", globalConfigDir)

	// Create global config and fragments
	createGlobalConfig(t, globalConfigDir)
	createGlobalFragments(t, globalConfigDir)

	// Create local fragments
	createLocalFragments(t, tmpDir)

	// Build the ctx binary
	ctxBinary := filepath.Join(tmpDir, "ctx")
	buildCmd := exec.Command("go", "build", "-o", ctxBinary, filepath.Join(originalWd, "cmd", "ctx"))
	buildCmd.Dir = originalWd

	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build ctx binary: %v", err)
	}

	return &integrationTestSetup{
		tmpDir:           tmpDir,
		ctxBinary:        ctxBinary,
		originalWd:       originalWd,
		oldXDGConfigHome: oldXDGConfigHome,
	}
}

func cleanupIntegrationTest(setup *integrationTestSetup) {
	_ = os.Chdir(setup.originalWd)
	if setup.oldXDGConfigHome != "" {
		_ = os.Setenv("XDG_CONFIG_HOME", setup.oldXDGConfigHome)
	} else {
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}
}

func createGlobalConfig(t *testing.T, globalConfigDir string) {
	configDir := filepath.Join(globalConfigDir, ".ctx")

	err := os.MkdirAll(configDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := `{
		"defaultTags": ["common", "typescript"],
		"outputFormats": {
			"test": "TEST.md"
		}
	}`
	configPath := filepath.Join(configDir, "config.json")

	err = os.WriteFile(configPath, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
}

func createGlobalFragments(t *testing.T, globalConfigDir string) {
	globalFragmentsDir := filepath.Join(globalConfigDir, ".ctx", "fragments")

	err := os.MkdirAll(globalFragmentsDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create global fragments directory: %v", err)
	}

	globalCommonContent := `---
ctx-tags: common, global
---

# Global Common Fragment
This is a global common fragment.
`

	err = os.WriteFile(filepath.Join(globalFragmentsDir, "common.md"), []byte(globalCommonContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create global common fragment: %v", err)
	}

	globalSpecificContent := `---
ctx-tags: typescript, global
---

# Global TypeScript Fragment
This is a global TypeScript fragment.
`

	err = os.WriteFile(filepath.Join(globalFragmentsDir, "typescript.md"), []byte(globalSpecificContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create global typescript fragment: %v", err)
	}
}

func createLocalFragments(t *testing.T, tmpDir string) {
	localFragmentsDir := filepath.Join(tmpDir, ".ctx", "fragments")

	err := os.MkdirAll(localFragmentsDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create local fragments directory: %v", err)
	}

	localCommonContent := `---
ctx-tags: common, local
---

# Local Common Fragment
This is a local common fragment that should override the global one.
`

	err = os.WriteFile(filepath.Join(localFragmentsDir, "common.md"), []byte(localCommonContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create local common fragment: %v", err)
	}

	localOnlyContent := `---
ctx-tags: local, testing
---

# Local Only Fragment
This fragment only exists locally.
`

	err = os.WriteFile(filepath.Join(localFragmentsDir, "local-only.md"), []byte(localOnlyContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create local-only fragment: %v", err)
	}
}

func TestLocalFragmentsIntegration_DefaultBehavior(t *testing.T) {
	setup := setupIntegrationTest(t)
	defer cleanupIntegrationTest(setup)

	// Run ctx build with default behavior (local overrides global)
	cmd := exec.Command(setup.ctxBinary, "build", "--non-interactive", "--stdout", "--tags", "common,typescript")
	cmd.Dir = setup.tmpDir

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("ctx build failed: %v", err)
	}

	outputStr := string(output)

	// Should contain local common fragment (not global)
	if !strings.Contains(outputStr, "Local Common Fragment") {
		t.Error("Expected output to contain local common fragment")
	}

	if strings.Contains(outputStr, "Global Common Fragment") {
		t.Error("Expected output to NOT contain global common fragment (should be overridden)")
	}

	// Should contain global typescript fragment
	if !strings.Contains(outputStr, "Global TypeScript Fragment") {
		t.Error("Expected output to contain global typescript fragment")
	}

	// Should NOT contain local-only fragment (not selected by tags)
	if strings.Contains(outputStr, "Local Only Fragment") {
		t.Error("Expected output to NOT contain local-only fragment (not selected by tags)")
	}
}

func TestLocalFragmentsIntegration_NoLocalOverride(t *testing.T) {
	setup := setupIntegrationTest(t)
	defer cleanupIntegrationTest(setup)

	// Run ctx build with --no-local-override flag, selecting only common tag
	cmd := exec.Command(setup.ctxBinary, "build", "--non-interactive", "--stdout", "--tags", "common", "--no-local-override")
	cmd.Dir = setup.tmpDir

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("ctx build with --no-local-override failed: %v", err)
	}

	outputStr := string(output)
	t.Logf("Output with --no-local-override: %s", outputStr)

	// Count occurrences of each fragment
	localCount := strings.Count(outputStr, "Local Common Fragment")
	globalCount := strings.Count(outputStr, "Global Common Fragment")

	// Should contain BOTH local and global common fragments (each appearing once)
	if localCount != 1 {
		t.Errorf("Expected local common fragment to appear once, got %d times", localCount)
	}

	if globalCount != 1 {
		t.Errorf("Expected global common fragment to appear once, got %d times", globalCount)
	}
}

func TestLocalFragmentsIntegration_LocalOnlyTags(t *testing.T) {
	setup := setupIntegrationTest(t)
	defer cleanupIntegrationTest(setup)

	// Run ctx build selecting only local tags
	cmd := exec.Command(setup.ctxBinary, "build", "--non-interactive", "--stdout", "--tags", "local,testing")
	cmd.Dir = setup.tmpDir

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("ctx build with local tags failed: %v", err)
	}

	outputStr := string(output)

	// Should contain both local fragments
	if !strings.Contains(outputStr, "Local Common Fragment") {
		t.Error("Expected output to contain local common fragment")
	}

	if !strings.Contains(outputStr, "Local Only Fragment") {
		t.Error("Expected output to contain local-only fragment")
	}

	// Should NOT contain global fragments
	if strings.Contains(outputStr, "Global Common Fragment") {
		t.Error("Expected output to NOT contain global common fragment")
	}

	if strings.Contains(outputStr, "Global TypeScript Fragment") {
		t.Error("Expected output to NOT contain global typescript fragment")
	}
}

func TestLocalFragmentsIntegration_HelpFlag(t *testing.T) {
	setup := setupIntegrationTest(t)
	defer cleanupIntegrationTest(setup)

	// Test that --no-local-override flag appears in help
	cmd := exec.Command(setup.ctxBinary, "build", "--help")
	cmd.Dir = setup.tmpDir

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("ctx build --help failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "--no-local-override") {
		t.Error("Expected help output to contain --no-local-override flag")
	}

	if !strings.Contains(outputStr, "include both local and global fragments") {
		t.Error("Expected help output to contain description of --no-local-override flag")
	}
}
