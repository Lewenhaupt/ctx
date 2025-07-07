package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFragment(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Fragment
	}{
		{
			name: "fragment with frontmatter",
			content: `---
ctx-tags: typescript, frontend, web
---

# TypeScript Guidelines

Some content here.`,
			expected: Fragment{
				Tags:    []string{"typescript", "frontend", "web"},
				Content: "\n# TypeScript Guidelines\n\nSome content here.",
			},
		},
		{
			name: "fragment without frontmatter",
			content: `# Simple Fragment

Just content, no tags.`,
			expected: Fragment{
				Tags:    []string{},
				Content: "# Simple Fragment\n\nJust content, no tags.",
			},
		},
		{
			name: "fragment with single tag",
			content: `---
ctx-tags: rust
---

# Rust Guidelines`,
			expected: Fragment{
				Tags:    []string{"rust"},
				Content: "\n# Rust Guidelines",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0o600)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Parse fragment
			fragment, err := ParseFragment(tmpFile)
			if err != nil {
				t.Fatalf("ParseFragment failed: %v", err)
			}

			// Check tags
			if len(fragment.Tags) != len(tt.expected.Tags) {
				t.Errorf("Expected %d tags, got %d", len(tt.expected.Tags), len(fragment.Tags))
			} else {
				for i, tag := range tt.expected.Tags {
					if i >= len(fragment.Tags) || fragment.Tags[i] != tag {
						t.Errorf("Expected tag %s at index %d, got %s", tag, i, fragment.Tags[i])
					}
				}
			}

			// Check content
			if fragment.Content != tt.expected.Content {
				t.Errorf("Expected content %q, got %q", tt.expected.Content, fragment.Content)
			}

			// Check path
			if fragment.Path != tmpFile {
				t.Errorf("Expected path %q, got %q", tmpFile, fragment.Path)
			}
		})
	}
}

func TestGetAllTags(t *testing.T) {
	fragments := []Fragment{
		{Tags: []string{"typescript", "frontend"}},
		{Tags: []string{"rust", "systems"}},
		{Tags: []string{"typescript", "web"}},
	}

	tags := GetAllTags(fragments)

	expected := map[string]bool{
		"typescript": true,
		"frontend":   true,
		"rust":       true,
		"systems":    true,
		"web":        true,
	}

	if len(tags) != len(expected) {
		t.Errorf("Expected %d unique tags, got %d", len(expected), len(tags))
	}

	for _, tag := range tags {
		if !expected[tag] {
			t.Errorf("Unexpected tag: %s", tag)
		}
	}
}

func TestFilterFragmentsByTags(t *testing.T) {
	fragments := []Fragment{
		{Path: "typescript.md", Tags: []string{"typescript", "frontend"}},
		{Path: "rust.md", Tags: []string{"rust", "systems"}},
		{Path: "general.md", Tags: []string{"general", "coding"}},
	}

	tests := []struct {
		name          string
		selectedTags  []string
		expectedPaths []string
	}{
		{
			name:          "single tag match",
			selectedTags:  []string{"typescript"},
			expectedPaths: []string{"typescript.md"},
		},
		{
			name:          "multiple tag match",
			selectedTags:  []string{"typescript", "rust"},
			expectedPaths: []string{"typescript.md", "rust.md"},
		},
		{
			name:          "no matches",
			selectedTags:  []string{"nonexistent"},
			expectedPaths: []string{},
		},
		{
			name:          "empty tags",
			selectedTags:  []string{},
			expectedPaths: []string{"typescript.md", "rust.md", "general.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterFragmentsByTags(fragments, tt.selectedTags)

			if len(filtered) != len(tt.expectedPaths) {
				t.Errorf("Expected %d fragments, got %d", len(tt.expectedPaths), len(filtered))
			}

			for i, expectedPath := range tt.expectedPaths {
				if i >= len(filtered) || filtered[i].Path != expectedPath {
					t.Errorf("Expected fragment path %s at index %d, got %s",
						expectedPath, i, filtered[i].Path)
				}
			}
		})
	}
}
