# Project Tasks

This file outlines the tasks required to build the `ctx` CLI tool.

## 1. Project Scaffolding & Setup
- [ ] Initialize Go module: `go mod init github.com/user/ctx`
- [ ] Set up basic project structure:
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
- [ ] Configure `flake.nix` to provide a development shell with Go and the required Charm.sh libraries.
- [ ] Configure `.envrc` to use the Nix direnv shell.

## 2. Configuration Management
- [ ] Define the JSON structure for the configuration file (`config.json`).
- [ ] Create a JSON schema (`config.schema.json`) for the configuration file.
- [ ] Implement logic to load config from `XDG_CONFIG_DIR/.ctx/config.json`.
- [ ] Implement the `--config-file` CLI argument to override the default config path.
- [ ] Add logic to read default tags from the configuration.

## 3. Core Logic
- [ ] Implement a file scanner to find all fragment files in `XDG_CONFIG_DIR/.ctx/fragments`.
- [ ] Implement a parser to read markdown files and extract `ctx-tags` from the frontmatter.
- [ ] Implement the file splicing logic to combine fragments based on selected tags.

## 4. CLI Implementation (`build` command)
- [ ] Set up the main command using `charm.sh/bubbletea`.
- [ ] Implement the interactive tag selection view:
    - [ ] Fetch and display all unique tags from fragments.
    - [ ] Use a multi-select component for tag selection.
    - [ ] Allow for pre-selection of tags via the `--tags` argument.
    - [ ] Add a de-selection step.
    - [ ] Add a final confirmation step before building.
- [ ] Implement the `--tags` flag for non-interactive use.
- [ ] After combining, output the result to `stdout`.

## 5. Replication & Reproducibility
- [ ] Implement the "command file" generation. This file will be created after a successful build and will contain a list of the fragment files used.
- [ ] Consider adding a `ctx build --from <command-file>` feature to allow rebuilding from this file.

## 6. Documentation & Best Practices
- [ ] Add comprehensive help text and documentation for all CLI commands, arguments, and configuration flags.
- [ ] Ensure all commits follow the `feat/chore/bug` prefix convention.
