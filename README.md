# ctx - Markdown Fragment Splicing Tool

A CLI tool for combining markdown fragments based on tags. Split your documentation, rules, or context files into multiple fragments and splice them together based on supplied tags.

## Features

- **Tag-based fragment selection**: Use `ctx-tags` in frontmatter to categorize fragments
- **Interactive TUI**: Select tags using an intuitive terminal interface
- **Non-interactive mode**: Automate builds with command-line flags
- **Configurable**: JSON configuration with schema validation
- **Multiple output formats**: Support for different AI tools (opencode, gemini, etc.)
- **Reproducible builds**: Generate command files for replication

## Installation

### Using Nix (Recommended)

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
go install github.com/user/ctx/cmd/ctx@latest
```

## Quick Start

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

## Configuration

Configuration is stored in `~/.config/.ctx/config.json` (or `$XDG_CONFIG_HOME/.ctx/config.json`).

### Example Configuration

```json
{
  "default_tags": ["general", "coding"],
  "output_formats": {
    "opencode": "AGENTS.md",
    "gemini": "GEMINI.md",
    "custom": "CUSTOM.md"
  },
  "fragments_dir": "/custom/path/to/fragments",
  "custom_settings": {
    "max_fragments": 50
  }
}
```

### Configuration Schema

The configuration follows the JSON schema defined in `config.schema.json`:

- `default_tags`: Array of tags to pre-select in interactive mode
- `output_formats`: Mapping of format names to output filenames
- `fragments_dir`: Custom path to fragments directory (optional)
- `custom_settings`: Additional settings for specific workflows

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
- Frontmatter is optional (fragments without tags are still valid)
- Only `.md` and `.markdown` files are processed
- Fragments are combined in the order they're found

## CLI Usage

```bash
ctx build [flags]

Flags:
  --tags strings         Comma-separated list of tags to include
  --non-interactive     Run in non-interactive mode
  --config-file string  Config file path (default: XDG_CONFIG_DIR/.ctx/config.json)
  -h, --help           Help for build
```

### Examples

```bash
# Interactive tag selection
ctx build

# Build with specific tags
ctx build --tags typescript,rust

# Non-interactive build with config file
ctx build --tags frontend --non-interactive --config-file ./my-config.json

# Build and save to file
ctx build --tags general,coding --non-interactive > output.md
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
go mod download

# Run tests
go test ./...

# Build
go build -o ctx ./cmd/ctx
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
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Follow commit message format: `feat/chore/bug: description`
6. Submit a pull request

## License

[Add your license here]