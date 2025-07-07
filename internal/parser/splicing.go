package parser

import (
	"fmt"
	"strings"
)

// SpliceFragments combines multiple fragments into a single output.
func SpliceFragments(fragments []Fragment) string {
	if len(fragments) == 0 {
		return ""
	}

	var result strings.Builder

	for i, fragment := range fragments {
		// Add a separator between fragments (except for the first one)
		if i > 0 {
			result.WriteString("\n\n")
		}

		// Add fragment content
		result.WriteString(fragment.Content)
	}

	return result.String()
}

// GenerateCommandFile creates a command file for replication.
func GenerateCommandFile(fragments []Fragment, selectedTags []string) string {
	var result strings.Builder

	result.WriteString("# ctx command file for replication\n")
	result.WriteString("# Generated automatically - do not edit manually\n\n")

	result.WriteString("## Selected Tags\n")

	for _, tag := range selectedTags {
		result.WriteString(fmt.Sprintf("- %s\n", tag))
	}

	result.WriteString("\n## Fragments Used\n")

	for _, fragment := range fragments {
		result.WriteString(fmt.Sprintf("- %s (tags: %s)\n",
			fragment.Path,
			strings.Join(fragment.Tags, ", ")))
	}

	result.WriteString(fmt.Sprintf("\n## Command\n```\nctx build --tags %s\n```\n",
		strings.Join(selectedTags, ",")))

	return result.String()
}
