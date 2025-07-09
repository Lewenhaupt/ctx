package main

import (
	"fmt"
	"os"

	"github.com/Lewenhaupt/ctx/internal/config"
	"github.com/Lewenhaupt/ctx/internal/parser"
	"github.com/Lewenhaupt/ctx/internal/tui"
	"github.com/spf13/cobra"
)

var (
	configFile      string
	tags            []string
	nonInteractive  bool
	outputFormats   []string
	outputFile      string
	stdout          bool
	noLocalOverride bool
)

var rootCmd = &cobra.Command{
	Use:   "ctx",
	Short: "A markdown splicing tool for combining fragments based on tags",
	Long: `ctx is a CLI tool for combining markdown fragments based on tags.
It allows users to split their files into multiple fragments with ctx-tags
and then splice them together based on supplied tags.`,
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and combine fragments based on tags",
	Long: `Build and combine markdown fragments based on selected tags.
The tool will scan the fragments directory, present available tags for selection,
and combine the matching fragments into a single output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := tui.BuildOptions{
			ConfigFile:      configFile,
			Tags:            tags,
			NonInteractive:  nonInteractive,
			OutputFormats:   outputFormats,
			OutputFile:      outputFile,
			Stdout:          stdout,
			NoLocalOverride: noLocalOverride,
		}
		return tui.RunBuild(&opts)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ctx configuration interactively",
	Long: `Initialize ctx configuration with an interactive questionnaire.
This command will guide you through setting up your ctx configuration,
creating the fragments directory, and optionally creating a sample fragment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := tui.InitOptions{
			ConfigFile: configFile,
		}
		return tui.RunInit(&opts)
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(ctx completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ ctx completion bash > /etc/bash_completion.d/ctx
  # macOS:
  $ ctx completion bash > $(brew --prefix)/etc/bash_completion.d/ctx

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ ctx completion zsh > "${fpath[1]}/_ctx"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ ctx completion fish | source

  # To load completions for each session, execute once:
  $ ctx completion fish > ~/.config/fish/completions/ctx.fish

PowerShell:

  PS> ctx completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> ctx completion powershell > ctx.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating completion: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "config file path (default: XDG_CONFIG_HOME/.ctx/config.json)")

	buildCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "comma-separated list of tags to include")
	buildCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "run in non-interactive mode")
	buildCmd.Flags().StringSliceVar(&outputFormats, "output-format", []string{}, "output format(s) to use (e.g., opencode, gemini, custom)")
	buildCmd.Flags().StringVar(&outputFile, "output-file", "", "output file path (overrides format-based naming)")
	buildCmd.Flags().BoolVar(&stdout, "stdout", false, "output to stdout instead of files")
	buildCmd.Flags().BoolVar(&noLocalOverride, "no-local-override", false, "include both local and global fragments even if they have the same name")

	// Add custom completion for tags flag
	if err := buildCmd.RegisterFlagCompletionFunc("tags", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAvailableTags(), cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering tags completion: %v\n", err)
	}

	// Add custom completion for output-format flag
	if err := buildCmd.RegisterFlagCompletionFunc("output-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"opencode", "gemini", "custom"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering output-format completion: %v\n", err)
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(completionCmd)
}

// getAvailableTags returns all available tags from fragments for completion.
func getAvailableTags() []string {
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return []string{}
	}

	fragmentsDir, err := config.GetFragmentsDir(cfg)
	if err != nil {
		return []string{}
	}

	globalFragments, err := parser.ScanFragments(fragmentsDir)
	if err != nil {
		return []string{}
	}

	localFragments, err := parser.ScanLocalFragments()
	if err != nil {
		return []string{}
	}

	fragments := parser.CombineFragments(globalFragments, localFragments, false)

	return parser.GetAllTags(fragments)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
