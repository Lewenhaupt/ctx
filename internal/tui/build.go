package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/user/ctx/internal/config"
	"github.com/user/ctx/internal/parser"
)

// BuildOptions represents the options for the build command
type BuildOptions struct {
	ConfigFile string
	Tags       []string
	NonInteractive bool
}

// RunBuild executes the build command with TUI
func RunBuild(opts BuildOptions) error {
	// Load configuration
	cfg, err := config.LoadConfig(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get fragments directory
	fragmentsDir, err := config.GetFragmentsDir(cfg)
	if err != nil {
		return fmt.Errorf("failed to get fragments directory: %w", err)
	}

	// Scan for fragments
	fragments, err := parser.ScanFragments(fragmentsDir)
	if err != nil {
		return fmt.Errorf("failed to scan fragments: %w", err)
	}

	if len(fragments) == 0 {
		return fmt.Errorf("no fragments found in %s", fragmentsDir)
	}

	// Get all available tags
	allTags := parser.GetAllTags(fragments)
	if len(allTags) == 0 {
		return fmt.Errorf("no tags found in fragments")
	}

	var selectedTags []string

	// If tags are provided via command line, use them
	if len(opts.Tags) > 0 {
		selectedTags = opts.Tags
	} else if opts.NonInteractive {
		// Use default tags from config in non-interactive mode
		selectedTags = cfg.DefaultTags
	} else {
		// Interactive tag selection
		selectedTags, err = selectTags(allTags, cfg.DefaultTags)
		if err != nil {
			return fmt.Errorf("tag selection failed: %w", err)
		}
	}

	// Filter fragments by selected tags
	filteredFragments := parser.FilterFragmentsByTags(fragments, selectedTags)
	if len(filteredFragments) == 0 {
		return fmt.Errorf("no fragments match the selected tags: %s", strings.Join(selectedTags, ", "))
	}

	// Show confirmation if interactive
	if !opts.NonInteractive {
		confirmed, err := confirmBuild(filteredFragments, selectedTags)
		if err != nil {
			return fmt.Errorf("confirmation failed: %w", err)
		}
		if !confirmed {
			fmt.Println("Build cancelled.")
			return nil
		}
	}

	// Splice fragments
	output := parser.SpliceFragments(filteredFragments)

	// Output to stdout
	fmt.Print(output)

	return nil
}

// selectTags presents an interactive multi-select for tag selection
func selectTags(allTags []string, defaultTags []string) ([]string, error) {
	var selectedTags []string

	// Create options for multi-select
	options := make([]huh.Option[string], len(allTags))
	for i, tag := range allTags {
		options[i] = huh.NewOption(tag, tag)
	}

	// Pre-select default tags
	defaultTagsMap := make(map[string]bool)
	for _, tag := range defaultTags {
		defaultTagsMap[tag] = true
	}

	var preSelected []string
	for _, tag := range allTags {
		if defaultTagsMap[tag] {
			preSelected = append(preSelected, tag)
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select tags to include:").
				Options(options...).
				Value(&selectedTags).
				Validate(func(val []string) error {
					if len(val) == 0 {
						return fmt.Errorf("at least one tag must be selected")
					}
					return nil
				}),
		),
	)

	// Set pre-selected values
	if len(preSelected) > 0 {
		selectedTags = preSelected
	}

	err := form.Run()
	if err != nil {
		return nil, err
	}

	return selectedTags, nil
}

// confirmBuild shows a confirmation dialog before building
func confirmBuild(fragments []parser.Fragment, selectedTags []string) (bool, error) {
	var confirmed bool

	// Create summary
	summary := fmt.Sprintf("Selected tags: %s\nFragments to include: %d\n\nFragments:\n", 
		strings.Join(selectedTags, ", "), len(fragments))
	
	for _, fragment := range fragments {
		summary += fmt.Sprintf("- %s (tags: %s)\n", fragment.Path, strings.Join(fragment.Tags, ", "))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Build Summary").
				Description(summary),
			huh.NewConfirm().
				Title("Proceed with build?").
				Value(&confirmed),
		),
	)

	err := form.Run()
	if err != nil {
		return false, err
	}

	return confirmed, nil
}