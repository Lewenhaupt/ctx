package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/ctx/internal/tui"
)

var (
	configFile     string
	tags           []string
	nonInteractive bool
	outputFormats  []string
	outputFile     string
	stdout         bool
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
			ConfigFile:     configFile,
			Tags:           tags,
			NonInteractive: nonInteractive,
			OutputFormats:  outputFormats,
			OutputFile:     outputFile,
			Stdout:         stdout,
		}
		return tui.RunBuild(&opts)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "config file path (default: XDG_CONFIG_DIR/.ctx/config.json)")

	buildCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "comma-separated list of tags to include")
	buildCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "run in non-interactive mode")
	buildCmd.Flags().StringSliceVar(&outputFormats, "output-format", []string{}, "output format(s) to use (e.g., opencode, gemini, custom)")
	buildCmd.Flags().StringVar(&outputFile, "output-file", "", "output file path (overrides format-based naming)")
	buildCmd.Flags().BoolVar(&stdout, "stdout", false, "output to stdout instead of files")

	rootCmd.AddCommand(buildCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
