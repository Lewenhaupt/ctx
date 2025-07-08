package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lewenhaupt/ctx/internal/config"
	"github.com/Lewenhaupt/ctx/internal/parser"
	"github.com/charmbracelet/huh"
)

// BuildOptions represents the options for the build command.
type BuildOptions struct {
	ConfigFile     string
	Tags           []string
	NonInteractive bool
	OutputFormats  []string
	OutputFile     string
	Stdout         bool
}

// RunBuild executes the build command with TUI.
func RunBuild(opts *BuildOptions) error {
	cfg, fragments, err := loadConfigAndFragments(opts.ConfigFile)
	if err != nil {
		return err
	}

	selectedTags, err := determineSelectedTags(opts, cfg, fragments)
	if err != nil {
		return err
	}

	filteredFragments := parser.FilterFragmentsByTags(fragments, selectedTags)
	if len(filteredFragments) == 0 {
		return fmt.Errorf("no fragments match the selected tags: %s", strings.Join(selectedTags, ", "))
	}

	selectedOutputFormats, outputFiles, err := determineOutputFormats(opts, cfg)
	if err != nil {
		return err
	}

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

	output := parser.SpliceFragments(filteredFragments)

	return handleOutput(opts, output, selectedOutputFormats, outputFiles, cfg)
}

func loadConfigAndFragments(configFile string) (*config.Config, []parser.Fragment, error) {
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	fragmentsDir, err := config.GetFragmentsDir(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get fragments directory: %w", err)
	}

	fragments, err := parser.ScanFragments(fragmentsDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan fragments: %w", err)
	}

	if len(fragments) == 0 {
		return nil, nil, fmt.Errorf("no fragments found in %s", fragmentsDir)
	}

	return cfg, fragments, nil
}

func determineSelectedTags(opts *BuildOptions, cfg *config.Config, fragments []parser.Fragment) ([]string, error) {
	allTags := parser.GetAllTags(fragments)
	if len(allTags) == 0 {
		return nil, fmt.Errorf("no tags found in fragments")
	}

	if len(opts.Tags) > 0 {
		return opts.Tags, nil
	}

	if opts.NonInteractive {
		return cfg.DefaultTags, nil
	}

	selectedTags, err := selectTags(allTags, cfg.DefaultTags)
	if err != nil {
		return nil, fmt.Errorf("tag selection failed: %w", err)
	}

	return selectedTags, nil
}

func determineOutputFormats(opts *BuildOptions, cfg *config.Config) (selectedFormats, outputFiles []string, err error) {
	if opts.Stdout {
		return []string{"stdout"}, nil, nil
	}

	if len(opts.OutputFormats) > 0 {
		return opts.OutputFormats, nil, nil
	}

	if opts.OutputFile != "" {
		return []string{"custom"}, []string{opts.OutputFile}, nil
	}

	if opts.NonInteractive {
		if len(cfg.OutputFormats) == 0 {
			return nil, nil, fmt.Errorf("no output formats configured and none specified")
		}

		var selectedOutputFormats []string

		for format := range cfg.OutputFormats {
			selectedOutputFormats = append(selectedOutputFormats, format)
		}

		return selectedOutputFormats, nil, nil
	}

	selectedOutputFormats, err := selectOutputFormats(cfg.OutputFormats)
	if err != nil {
		return nil, nil, fmt.Errorf("output format selection failed: %w", err)
	}

	return selectedOutputFormats, nil, nil
}

func handleOutput(opts *BuildOptions, output string, selectedOutputFormats, outputFiles []string, cfg *config.Config) error {
	if opts.Stdout {
		fmt.Print(output)
		return nil
	}

	err := writeOutputFiles(output, selectedOutputFormats, outputFiles, cfg)
	if err != nil {
		return fmt.Errorf("failed to write output files: %w", err)
	}

	return nil
}

// selectTags presents an interactive multi-select for tag selection.
func selectTags(allTags, defaultTags []string) ([]string, error) {
	// Pre-select default tags that exist in allTags
	defaultTagsMap := make(map[string]bool)
	for _, tag := range defaultTags {
		defaultTagsMap[tag] = true
	}

	var selectedTags []string
	for _, tag := range allTags {
		if defaultTagsMap[tag] {
			selectedTags = append(selectedTags, tag)
		}
	}

	// Create options for multi-select
	options := make([]huh.Option[string], len(allTags))
	for i, tag := range allTags {
		options[i] = huh.NewOption(tag, tag)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select tags to include:").
				Description("Use space to toggle selection, enter to confirm").
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

	err := form.Run()
	if err != nil {
		return nil, err
	}

	return selectedTags, nil
}

// selectOutputFormats presents an interactive multi-select for output format selection.
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

// confirmBuild shows a confirmation dialog before building.
func confirmBuild(fragments []parser.Fragment, selectedTags, outputFormats []string) (bool, error) {
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

// writeOutputFiles writes the output to the specified files based on formats.
func writeOutputFiles(output string, formats, customFiles []string, cfg *config.Config) error {
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
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		// Write file
		if err := os.WriteFile(filename, []byte(output), 0o600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}

		fmt.Printf("Output written to: %s\n", filename)
	}

	return nil
}
