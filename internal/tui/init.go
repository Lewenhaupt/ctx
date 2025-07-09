package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Lewenhaupt/ctx/internal/config"
	"github.com/charmbracelet/huh"
)

// InitOptions represents the options for the init command.
type InitOptions struct {
	ConfigFile string
}

// InitAnswers holds the user's responses to the init questionnaire.
type InitAnswers struct {
	AddOutputFormats bool
	FragmentsDir     string
	CreateSample     bool
}

// RunInit executes the init command with interactive questionnaire.
func RunInit(opts *InitOptions) error {
	// Check if config already exists
	configPath := opts.ConfigFile
	if configPath == "" {
		configDir, err := config.GetConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}

		configPath = filepath.Join(configDir, "config.json")
	}

	if _, err := os.Stat(configPath); err == nil {
		var overwrite bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Configuration file already exists").
					Description(fmt.Sprintf("A configuration file already exists at %s. Do you want to overwrite it?", configPath)).
					Value(&overwrite),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get overwrite confirmation: %w", err)
		}

		if !overwrite {
			fmt.Println("Init cancelled.")
			return nil
		}
	}

	// Run interactive questionnaire
	answers, err := runQuestionnaire()
	if err != nil {
		return fmt.Errorf("questionnaire failed: %w", err)
	}

	// Generate configuration
	cfg, err := generateConfig(answers)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Save configuration
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration saved to: %s\n", configPath)

	// Create fragments directory
	fragmentsDir, err := config.GetFragmentsDir(cfg)
	if err != nil {
		return fmt.Errorf("failed to get fragments directory: %w", err)
	}

	if err := os.MkdirAll(fragmentsDir, 0o750); err != nil {
		return fmt.Errorf("failed to create fragments directory: %w", err)
	}

	fmt.Printf("Fragments directory created: %s\n", fragmentsDir)

	// Create sample fragment if requested
	if answers.CreateSample {
		if err := createSampleFragment(fragmentsDir); err != nil {
			return fmt.Errorf("failed to create sample fragment: %w", err)
		}
	}

	fmt.Println("\nSetup complete! You can now run 'ctx build' to start using the tool.")

	return nil
}

// runQuestionnaire presents the interactive questionnaire to the user.
func runQuestionnaire() (*InitAnswers, error) {
	answers := &InitAnswers{}

	// Get default config to show default output formats
	defaultCfg := config.DefaultConfig()

	var outputFormatsDesc string
	for format, filename := range defaultCfg.OutputFormats {
		outputFormatsDesc += fmt.Sprintf("- %s: %s\n", format, filename)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome to ctx init!").
				Description("This will help you set up your ctx configuration."),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("Default Output Formats").
				Description("The following output formats are available by default:\n"+outputFormatsDesc),
			huh.NewConfirm().
				Title("Do you want to add additional output formats?").
				Description("You can add custom output formats beyond the defaults").
				Value(&answers.AddOutputFormats),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Where would you like to store your fragments?").
				Description("This is where your markdown fragments will be stored").
				Placeholder("./fragments").
				Value(&answers.FragmentsDir).
				Validate(func(val string) error {
					if val == "" {
						return nil // Allow empty for default
					}
					// Basic path validation
					if filepath.IsAbs(val) {
						return nil
					}

					if _, err := filepath.Abs(val); err != nil {
						return fmt.Errorf("invalid path: %w", err)
					}

					return nil
				}),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to create a sample fragment to start?").
				Description("This will create a hello-world example fragment").
				Value(&answers.CreateSample),
		),
	)
	if err := form.Run(); err != nil {
		return nil, err
	}

	return answers, nil
}

// generateConfig creates a configuration based on user answers.
func generateConfig(answers *InitAnswers) (*config.Config, error) {
	cfg := config.DefaultConfig()

	// Set fragments directory if provided
	if answers.FragmentsDir != "" {
		// Convert relative path to absolute
		absPath, err := filepath.Abs(answers.FragmentsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to convert fragments directory to absolute path: %w", err)
		}

		cfg.FragmentsDir = absPath
	}

	// Handle additional output formats if user requested them
	if answers.AddOutputFormats {
		customFormats, err := promptForCustomFormats()
		if err != nil {
			return nil, fmt.Errorf("failed to get custom formats: %w", err)
		}

		// Add custom formats to the config
		for name, filename := range customFormats {
			cfg.OutputFormats[name] = filename
		}
	}

	return cfg, nil
}

// createSampleFragment creates a hello-world sample fragment.
func createSampleFragment(fragmentsDir string) error {
	sampleContent := `ctx-tags: hello, world, sample

# Hello World

This is a sample fragment created by ctx init.

You can edit this file and add more fragments to get started with ctx.

## Usage

Run 'ctx build' to combine fragments based on tags.
`

	samplePath := filepath.Join(fragmentsDir, "hello-world.md")
	if err := os.WriteFile(samplePath, []byte(sampleContent), 0o600); err != nil {
		return fmt.Errorf("failed to write sample fragment: %w", err)
	}

	fmt.Printf("Sample fragment created: %s\n", samplePath)

	return nil
}

// promptForCustomFormats prompts the user to add custom output formats.
func promptForCustomFormats() (map[string]string, error) {
	customFormats := make(map[string]string)

	addMore := true

	for addMore {
		var formatName, fileName string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Format name").
					Description("Enter a name for the custom output format (e.g., 'claude', 'custom')").
					Value(&formatName).
					Validate(func(val string) error {
						if val == "" {
							return fmt.Errorf("format name cannot be empty")
						}

						if _, exists := customFormats[val]; exists {
							return fmt.Errorf("format name already exists")
						}

						return nil
					}),
				huh.NewInput().
					Title("File name").
					Description("Enter the output file name for this format (e.g., 'CLAUDE.md', 'output.txt')").
					Value(&fileName).
					Validate(func(val string) error {
						if val == "" {
							return fmt.Errorf("file name cannot be empty")
						}

						return nil
					}),
			),
		)
		if err := form.Run(); err != nil {
			return nil, err
		}

		customFormats[formatName] = fileName
		// Ask if user wants to add more formats
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another custom output format?").
					Value(&addMore),
			),
		)

		if err := form.Run(); err != nil {
			return nil, err
		}
	}

	return customFormats, nil
}
