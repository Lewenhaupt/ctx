package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanLocalFragments(t *testing.T) {
	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWd)
	}()

	// Create temporary directory and change to it
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test case 1: No local fragments directory
	fragments, err := ScanLocalFragments()
	if err != nil {
		t.Fatalf("ScanLocalFragments failed when no directory exists: %v", err)
	}
	if len(fragments) != 0 {
		t.Errorf("Expected 0 fragments when no directory exists, got %d", len(fragments))
	}

	// Test case 2: Create local fragments directory with files
	localFragmentsDir := filepath.Join(tmpDir, ".ctx", "fragments")
	err = os.MkdirAll(localFragmentsDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create local fragments directory: %v", err)
	}

	// Create test fragments
	fragment1Content := `---
ctx-tags: local, typescript
---

# Local TypeScript Fragment
This is a local fragment.`

	fragment2Content := `---
ctx-tags: local, testing
---

# Local Testing Fragment
This is another local fragment.`

	err = os.WriteFile(filepath.Join(localFragmentsDir, "local-ts.md"), []byte(fragment1Content), 0o600)
	if err != nil {
		t.Fatalf("Failed to create local fragment file: %v", err)
	}

	err = os.WriteFile(filepath.Join(localFragmentsDir, "local-test.md"), []byte(fragment2Content), 0o600)
	if err != nil {
		t.Fatalf("Failed to create local fragment file: %v", err)
	}

	// Scan local fragments
	fragments, err = ScanLocalFragments()
	if err != nil {
		t.Fatalf("ScanLocalFragments failed: %v", err)
	}

	if len(fragments) != 2 {
		t.Errorf("Expected 2 local fragments, got %d", len(fragments))
	}

	// Verify fragments contain expected tags
	foundTags := make(map[string]bool)
	for _, fragment := range fragments {
		for _, tag := range fragment.Tags {
			foundTags[tag] = true
		}
	}

	expectedTags := []string{"local", "typescript", "testing"}
	for _, expectedTag := range expectedTags {
		if !foundTags[expectedTag] {
			t.Errorf("Expected to find tag %s in local fragments", expectedTag)
		}
	}
}

func TestCombineFragments(t *testing.T) {
	globalFragments := []Fragment{
		{Path: "/global/common.md", Tags: []string{"global", "common"}, Content: "Global common content"},
		{Path: "/global/specific.md", Tags: []string{"global", "specific"}, Content: "Global specific content"},
	}

	localFragments := []Fragment{
		{Path: "/local/.ctx/fragments/common.md", Tags: []string{"local", "common"}, Content: "Local common content"},
		{Path: "/local/.ctx/fragments/local-only.md", Tags: []string{"local", "only"}, Content: "Local only content"},
	}

	t.Run("default behavior - local overrides global", func(t *testing.T) {
		combined := CombineFragments(globalFragments, localFragments, false)

		// Should have 3 fragments: local common (overriding global), global specific, local only
		if len(combined) != 3 {
			t.Errorf("Expected 3 combined fragments, got %d", len(combined))
		}

		// Check that local common.md overrode global common.md
		var commonFragment *Fragment
		for i, fragment := range combined {
			if filepath.Base(fragment.Path) == "common.md" {
				commonFragment = &combined[i]
				break
			}
		}

		if commonFragment == nil {
			t.Error("Expected to find common.md fragment")
		} else if commonFragment.Content != "Local common content" {
			t.Errorf("Expected local common content to override global, got: %s", commonFragment.Content)
		}

		// Verify all expected fragments are present
		expectedPaths := map[string]bool{
			"common.md":     true,
			"specific.md":   true,
			"local-only.md": true,
		}

		for _, fragment := range combined {
			basename := filepath.Base(fragment.Path)
			if !expectedPaths[basename] {
				t.Errorf("Unexpected fragment: %s", basename)
			}
			delete(expectedPaths, basename)
		}

		if len(expectedPaths) > 0 {
			t.Errorf("Missing expected fragments: %v", expectedPaths)
		}
	})

	t.Run("no-local-override - include both local and global", func(t *testing.T) {
		combined := CombineFragments(globalFragments, localFragments, true)

		// Should have 4 fragments: all global + all local
		if len(combined) != 4 {
			t.Errorf("Expected 4 combined fragments, got %d", len(combined))
		}

		// Check that both common.md files are present
		commonCount := 0
		for _, fragment := range combined {
			if filepath.Base(fragment.Path) == "common.md" {
				commonCount++
			}
		}

		if commonCount != 2 {
			t.Errorf("Expected 2 common.md fragments when no-local-override is true, got %d", commonCount)
		}
	})

	t.Run("empty fragments", func(t *testing.T) {
		// Test with empty global fragments
		combined := CombineFragments([]Fragment{}, localFragments, false)
		if len(combined) != len(localFragments) {
			t.Errorf("Expected %d fragments with empty global, got %d", len(localFragments), len(combined))
		}

		// Test with empty local fragments
		combined = CombineFragments(globalFragments, []Fragment{}, false)
		if len(combined) != len(globalFragments) {
			t.Errorf("Expected %d fragments with empty local, got %d", len(globalFragments), len(combined))
		}

		// Test with both empty
		combined = CombineFragments([]Fragment{}, []Fragment{}, false)
		if len(combined) != 0 {
			t.Errorf("Expected 0 fragments with both empty, got %d", len(combined))
		}
	})
}
