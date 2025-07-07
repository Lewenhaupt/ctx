package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/user/ctx/internal/config"
	"github.com/user/ctx/internal/parser"
)

// BuildOptions represents the options for the build command
type BuildOptions struct {
	ConfigFile     string
	Tags           []string
	NonInteractive bool
	OutputFormats  []string
	OutputFile     string
	Stdout         bool
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

	var selectedOutputFormats []string
	var outputFiles []string

	// Handle output format selection
	if opts.Stdout {
		// Output to stdout - no format selection needed
		selectedOutputFormats = []string{"stdout"}
	} else if len(opts.OutputFormats) > 0 {
		// Use provided output formats
		selectedOutputFormats = opts.OutputFormats
	} else if opts.OutputFile != "" {
		// Use custom output file
		selectedOutputFormats = []string{"custom"}
		outputFiles = []string{opts.OutputFile}
	} else if opts.NonInteractive {
		// Use default output formats from config in non-interactive mode
		if len(cfg.OutputFormats) == 0 {
			return fmt.Errorf("no output formats configured and none specified")
		}
		for format := range cfg.OutputFormats {
			selectedOutputFormats = append(selectedOutputFormats, format)
		}
	} else {
		// Interactive output format selection
		selectedOutputFormats, err = selectOutputFormats(cfg.OutputFormats)
		if err != nil {
			return fmt.Errorf("output format selection failed: %w", err)
		}
	}

	// Show confirmation if interactive
	if !opts.NonInteractive {
		confirmed, err := confirmBuild(filteredFragments, selectedTags, selectedOutputFormats)
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

	// Handle output
	if opts.Stdout {
		// Output to stdout
		fmt.Print(output)
	} else {
		// Write to files
		err = writeOutputFiles(output, selectedOutputFormats, outputFiles, cfg)
		if err != nil {
			return fmt.Errorf("failed to write output files: %w", err)
		}
	}

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

// selectOutputFormats presents an interactive multi-select for output format selection
func selectOutputFormats(availableFormats map[string]string) ([]string, error) {
	var selectedFormats []string

	// Create options for multi-select
	options := make([]huh.Option[string], 0, len(availableFormats)+1)

	// Add configured formats
	for format := range availableFormats {
		options = append(options, huh.NewOption(format, format))
	}

	// Add stdout option
	options = append(options, huh.NewOption("stdout", "stdout"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select output format(s):").
				Options(options...).
				Value(&selectedFormats).
				Validate(func(val []string) error {
					if len(val) == 0 {
						return fmt.Errorf("at least one output format must be selected")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, err
	}

	return selectedFormats, nil
}

// confirmBuild shows a confirmation dialog before building
func confirmBuild(fragments []parser.Fragment, selectedTags []string, outputFormats []string) (bool, error) {
	var confirmed bool

	// Create summary
	summary := fmt.Sprintf("Selected tags: %s\nOutput formats: %s\nFragments to include: %d\n\nFragments:\n",
		strings.Join(selectedTags, ", "), strings.Join(outputFormats, ", "), len(fragments))

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

// writeOutputFiles writes the output to the specified files based on formats
func writeOutputFiles(output string, formats []string, customFiles []string, cfg *config.Config) error {
	for i, format := range formats {
		var filename string

		if format == "stdout" {
			// Skip stdout in file writing
			continue
		} else if format == "custom" && i < len(customFiles) {
			filename = customFiles[i]
		} else if outputFile, exists := cfg.OutputFormats[format]; exists {
			filename = outputFile
		} else {
			return fmt.Errorf("unknown output format: %s", format)
		}

		// Ensure directory exists
		dir := filepath.Dir(filename)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		// Write file
		if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}

		fmt.Printf("Output written to: %s\n", filename)
	}

	return nil
}
