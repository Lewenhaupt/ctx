package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Fragment represents a markdown fragment with its metadata
type Fragment struct {
	Path    string   `json:"path"`
	Tags    []string `json:"tags"`
	Content string   `json:"content"`
}

// ScanFragments scans the fragments directory and returns all found fragments
func ScanFragments(fragmentsDir string) ([]Fragment, error) {
	var fragments []Fragment

	if _, err := os.Stat(fragmentsDir); os.IsNotExist(err) {
		return fragments, nil // Return empty slice if directory doesn't exist
	}

	err := filepath.Walk(fragmentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process markdown files
		if !info.IsDir() && (strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")) {
			fragment, err := ParseFragment(path)
			if err != nil {
				return fmt.Errorf("failed to parse fragment %s: %w", path, err)
			}
			fragments = append(fragments, *fragment)
		}

		return nil
	})

	return fragments, err
}

// ParseFragment parses a single markdown file and extracts ctx-tags and content
func ParseFragment(filePath string) (*Fragment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tags []string
	var contentLines []string
	var inFrontmatter bool
	var frontmatterProcessed bool

	// Regex to match ctx-tags line
	ctxTagsRegex := regexp.MustCompile(`^ctx-tags:\s*(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for frontmatter boundaries
		if strings.TrimSpace(line) == "---" {
			if !frontmatterProcessed {
				inFrontmatter = !inFrontmatter
				if !inFrontmatter {
					frontmatterProcessed = true
				}
				continue
			}
		}

		// If we're in frontmatter, look for ctx-tags
		if inFrontmatter {
			if matches := ctxTagsRegex.FindStringSubmatch(line); matches != nil {
				tagsStr := strings.TrimSpace(matches[1])
				// Split by comma and clean up each tag
				for _, tag := range strings.Split(tagsStr, ",") {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						tags = append(tags, tag)
					}
				}
			}
		} else if frontmatterProcessed {
			// After frontmatter, collect all content
			contentLines = append(contentLines, line)
		} else {
			// No frontmatter detected, treat as content
			contentLines = append(contentLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	content := strings.Join(contentLines, "\n")

	return &Fragment{
		Path:    filePath,
		Tags:    tags,
		Content: content,
	}, nil
}

// GetAllTags extracts all unique tags from a slice of fragments
func GetAllTags(fragments []Fragment) []string {
	tagSet := make(map[string]bool)
	for _, fragment := range fragments {
		for _, tag := range fragment.Tags {
			tagSet[tag] = true
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags
}

// FilterFragmentsByTags returns fragments that contain any of the specified tags
func FilterFragmentsByTags(fragments []Fragment, selectedTags []string) []Fragment {
	if len(selectedTags) == 0 {
		return fragments
	}

	tagSet := make(map[string]bool)
	for _, tag := range selectedTags {
		tagSet[tag] = true
	}

	var filtered []Fragment
	for _, fragment := range fragments {
		for _, tag := range fragment.Tags {
			if tagSet[tag] {
				filtered = append(filtered, fragment)
				break
			}
		}
	}

	return filtered
}