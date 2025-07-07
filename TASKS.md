# Project Tasks

This file outlines the tasks required to build the `ctx` CLI tool.

## 1. Project Scaffolding & Setup ✅
- [x] Initialize Go module: `go mod init github.com/user/ctx`
- [x] Set up basic project structure:
  ```
  .
  ├── cmd/
  │   └── ctx/
  │       └── main.go
  ├── internal/
  │   ├── config/
  │   ├── tui/
  │   └── parser/
  └── go.mod
  ```
- [x] Configure `flake.nix` to provide a development shell with Go and the required Charm.sh libraries.
- [x] Configure `.envrc` to use the Nix direnv shell.

## 2. Configuration Management ✅
- [x] Define the JSON structure for the configuration file (`config.json`).
- [x] Create a JSON schema (`config.schema.json`) for the configuration file.
- [x] Implement logic to load config from `XDG_CONFIG_HOME/.ctx/config.json`.
- [x] Implement the `--config-file` CLI argument to override the default config path.
- [x] Add logic to read default tags from the configuration.

## 3. Core Logic ✅
- [x] Implement a file scanner to find all fragment files in `XDG_CONFIG_HOME/.ctx/fragments`.
- [x] Implement a parser to read markdown files and extract `ctx-tags` from the frontmatter.
- [x] Implement the file splicing logic to combine fragments based on selected tags.

## 4. CLI Implementation (`build` command) ✅
- [x] Set up the main command using `charm.sh/bubbletea`.
- [x] Implement the interactive tag selection view:
    - [x] Fetch and display all unique tags from fragments.
    - [x] Use a multi-select component for tag selection.
    - [x] Allow for pre-selection of tags via the `--tags` argument.
    - [x] Add a de-selection step.
    - [x] Add a final confirmation step before building.
- [x] Implement the `--tags` flag for non-interactive use.
- [x] After combining, output the result to `stdout`.

## 5. Replication & Reproducibility ✅
- [x] Implement the "command file" generation. This file will be created after a successful build and will contain a list of the fragment files used.
- [ ] Consider adding a `ctx build --from <command-file>` feature to allow rebuilding from this file.

## 6. Documentation & Best Practices ✅
- [x] Add comprehensive help text and documentation for all CLI commands, arguments, and configuration flags.
- [x] Ensure all commits follow the `feat/chore/bug` prefix convention.
- [x] Add unit tests for all core functionality
- [x] Add integration tests for end-to-end functionality
- [x] Create comprehensive README documentation

## Status: ✅ COMPLETED

The `ctx` CLI tool has been successfully implemented with all core features:

- ✅ Interactive and non-interactive modes
- ✅ Tag-based fragment selection
- ✅ Configuration management with JSON schema
- ✅ Comprehensive test coverage
- ✅ Full documentation
- ✅ Nix development environment

### Future Enhancements (Optional)
- [ ] Add `ctx build --from <command-file>` feature for reproducible builds
- [ ] Add support for nested fragment directories
- [ ] Implement fragment validation and linting
- [ ] Add support for custom output templates
- [ ] Create shell completion scripts
