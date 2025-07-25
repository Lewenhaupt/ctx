# ctx - Markdown Fragment Splicing Tool

A CLI tool for combining markdown fragments based on tags. Split your documentation, rules, or context files into multiple fragments and splice them together based on supplied tags.

## Features

- **Tag-based fragment selection**: Use `ctx-tags` in frontmatter to categorize fragments
- **Interactive TUI**: Select tags using an intuitive terminal interface
- **Non-interactive mode**: Automate builds with command-line flags
- **Configurable**: JSON configuration with schema validation
- **Multiple output formats**: Support for different AI tools (opencode, gemini, etc.)
- **Project-specific fragments**: Local `.ctx/fragments` directory support with override logic
- **Reproducible builds**: Generate command files for replication

## Installation

### Using Nix (Recommended)

#### Option 1: Direct Installation from GitHub

```bash
# Install directly from the repository
nix profile install github:Lewenhaupt/ctx

# Or run without installing
nix run github:Lewenhaupt/ctx -- build --help

# For faster builds, you can also use the binary cache (if available)
# nix profile install github:Lewenhaupt/ctx --extra-substituters https://cache.nixos.org
```

#### Option 2: Using Nix Flakes in your project

Add to your `flake.nix`:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    ctx.url = "github:Lewenhaupt/ctx";
  };

  outputs = { self, nixpkgs, ctx, ... }:
    let
      system = "x86_64-linux"; # or your system
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.${system}.default = pkgs.mkShell {
        packages = [
          ctx.packages.${system}.default
          # other packages...
        ];
      };
    };
}
```

Then use in your project:

```bash
# Enter development shell with ctx available
nix develop

# Or run directly
nix run .#ctx -- build --help
```

#### Option 3: NixOS System Installation

Add to your NixOS configuration:

```nix
# configuration.nix or flake.nix
{
  inputs.ctx.url = "github:Lewenhaupt/ctx";
  
  # In your system configuration:
  environment.systemPackages = [
    inputs.ctx.packages.${system}.default
  ];
}
```

#### Option 4: Development Setup

```bash
# Clone the repository
git clone <repository-url>
cd ctx

# Enter development shell
direnv allow
# or
nix develop

# Build the binary
go build -o ctx ./cmd/ctx
```

### Using Go

```bash
go install github.com/Lewenhaupt/ctx/cmd/ctx@latest
```

## Getting Started

### Initialize Configuration

The easiest way to get started is to use the interactive init command to scaffold your configuration:

```bash
ctx init
```

This will guide you through:
- Setting up your configuration file (`~/.config/.ctx/config.json`)
- Choosing output formats (opencode, gemini, or custom formats)
- Configuring your fragments directory location
- Optionally creating a sample fragment to get started

The init command will:
1. Show you the default output formats available
2. Allow you to add custom output formats if needed
3. Let you specify where to store your fragments (defaults to `~/.config/.ctx/fragments`)
4. Optionally create a hello-world sample fragment
5. Create all necessary directories and configuration files

### Manual Setup (Alternative)

If you prefer to set up manually:

1. **Create fragments directory**:
   ```bash
   mkdir -p ~/.config/.ctx/fragments
   ```

2. **Create a fragment** (`~/.config/.ctx/fragments/typescript.md`):
   ```markdown
   ---
   ctx-tags: typescript, frontend, web
   ---

   # TypeScript Guidelines

   ## Type Safety
   - Always use strict mode
   - Prefer interfaces over types for object shapes
   ```

3. **Build fragments**:
   ```bash
   # Interactive mode
   ctx build

   # Non-interactive mode
   ctx build --tags typescript,frontend

   # Output to file
   ctx build --tags typescript --non-interactive > AGENTS.md
   ```

## Quick Start

After running `ctx init`, you can immediately start building:

```bash
# Interactive mode - select tags and output formats
ctx build

# Non-interactive mode with specific tags
ctx build --tags hello,world --non-interactive
```

## Configuration

Configuration is stored in `~/.config/.ctx/config.json` (or `$XDG_CONFIG_HOME/.ctx/config.json`).

### Example Configuration

```json
{
  "defaultTags": ["general", "coding"],
  "outputFormats": {
    "opencode": "AGENTS.md",
    "gemini": "GEMINI.md",
    "custom": "CUSTOM.md"
  },
  "fragmentsDir": "/custom/path/to/fragments",
  "customSettings": {
    "max_fragments": 50
  }
}
```

### Configuration Schema

The configuration follows the JSON schema defined in `config.schema.json`:

- `defaultTags`: Array of tags to pre-select in interactive mode
- `outputFormats`: Mapping of format names to output filenames
- `fragmentsDir`: Custom path to fragments directory (optional)
- `customSettings`: Additional settings for specific workflows

## Fragment Format

Fragments are markdown files with optional frontmatter containing `ctx-tags`:

```markdown
---
ctx-tags: tag1, tag2, tag3
---

# Fragment Content

Your markdown content here.
```

### Rules

- Tags are comma-separated in the `ctx-tags` field
- Frontmatter is optional 
- Only `.md` and `.markdown` files are processed
- Fragments are combined in the order they're found

## Project-Specific Fragments

In addition to global fragments stored in your configuration directory, `ctx` supports project-specific fragments stored in a local `.ctx/fragments` directory within your project.

### Local Fragment Directory

Create a `.ctx/fragments` directory in your project root:

```bash
mkdir -p .ctx/fragments
```

Add project-specific fragments:

```bash
# .ctx/fragments/api-guidelines.md
---
ctx-tags: api, guidelines, backend
---

# API Guidelines

These guidelines apply specifically to this project's API.
```

### Override Behavior

By default, local fragments override global fragments with the same filename:

- **Global**: `~/.config/.ctx/fragments/common.md`
- **Local**: `./.ctx/fragments/common.md`
- **Result**: Local `common.md` is used, global is ignored

### Including Both Local and Global

Use the `--no-local-override` flag to include both local and global fragments:

```bash
# Include both local and global fragments with the same name
ctx build --no-local-override --tags common
```

### Examples

```bash
# Build using both global and local fragments (default override behavior)
ctx build --tags api,backend

# Build including both local and global fragments with same names
ctx build --tags common --no-local-override

# Build with tags that might exist in both global and local fragments
ctx build --tags guidelines,typescript
```

### Use Cases

- **Project documentation**: Store project-specific context in `.ctx/fragments`
- **Team overrides**: Override global rules with project-specific ones
- **Environment-specific configs**: Different fragments for different projects
- **Version control**: Commit local fragments with your project

## CLI Commands

### Initialize Configuration

```bash
ctx init [flags]

Flags:
  --config-file string   Config file path (default: XDG_CONFIG_HOME/.ctx/config.json)
  -h, --help            Help for init
```

### Build Fragments

```bash
ctx build [flags]

Flags:
  --tags strings              Comma-separated list of tags to include
  --non-interactive          Run in non-interactive mode
  --output-format strings    Output format(s) to use (e.g., opencode, gemini, custom)
  --output-file string       Output file path (overrides format-based naming)
  --stdout                   Output to stdout instead of files
  --no-local-override        Include both local and global fragments even if they have the same name
  --config-file string       Config file path (default: XDG_CONFIG_HOME/.ctx/config.json)
  -h, --help                Help for build
```

### Examples

```bash
# Interactive mode (select tags and output formats)
ctx build

# Build with specific tags and output to stdout
ctx build --tags typescript,rust --stdout

# Non-interactive build with specific output format
ctx build --tags frontend --non-interactive --output-format opencode

# Build with multiple output formats
ctx build --tags general,coding --output-format opencode,gemini --non-interactive

# Build to custom file
ctx build --tags typescript --output-file custom-output.md --non-interactive

# Pipe output to another command
ctx build --tags general --stdout --non-interactive | grep "function"
```

## Development

### Prerequisites

- Go 1.21+
- Nix (optional, for development environment)

### Setup

```bash
# Clone repository
git clone <repository-url>
cd ctx

# Setup development environment
direnv allow  # if using Nix
# or
nix develop  # if using Nix without direnv
# or
go mod download  # if using Go directly

# Setup git hooks for commit message validation
./scripts/setup-git-hooks.sh

# Run tests
go test ./...

# Build using Go
go build -o ctx ./cmd/ctx

# Or build using Nix
nix build .#default
```

### Running Locally During Development

```bash
# Run directly with go run
go run ./cmd/ctx build

# Run with specific flags and output to stdout
go run ./cmd/ctx build --tags typescript,rust --stdout

# Run non-interactively with specific output format
go run ./cmd/ctx build --tags general --non-interactive --output-format opencode

# Use direnv to ensure correct environment
direnv exec . go run ./cmd/ctx build

# Build and run the binary
go build -o ctx ./cmd/ctx
./ctx build --stdout --non-interactive --tags general
```

### Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests
go test .

# All tests with coverage
go test -cover ./...
```

### Project Structure

```
.
├── cmd/ctx/           # Main CLI application
├── internal/
│   ├── config/        # Configuration management
│   ├── parser/        # Fragment parsing and splicing
│   └── tui/          # Terminal UI components
├── config.schema.json # JSON schema for configuration
├── flake.nix         # Nix development environment
└── integration_test.go # Integration tests
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Set up git hooks: `./scripts/setup-git-hooks.sh`
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Follow commit message format: `type: description` (enforced by git hooks)
   - Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
   - Example: `feat: add new fragment parsing feature`
7. Submit a pull request

## License

[Add your license here]# Test commit hook
test change
