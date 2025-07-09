# ctx rules/markdown splicing tool mono-repo

This project provides a cli tool for combining rules/markdown files/context/etc. for AI tools. It allows users to split their file(s) into multiple fragments. These markdown fragments should contain a ctx-tags entry at the top which is a list of tags. The tool will then allow the user to splice these files together based on supplied tags.
Example command: ```ctx build --tags typescript, rust```

## Tech
- Language: GO
- Frameworks: charm.sh (bubble tea, huh, lip gloss, bubbles, log)
- git
- Nix (for building and managing dev-shell, uses direnv)
- `direnv exec` for making sure we always execute our commands in the correct dev-shell.

## Best practices
- Always make commits after making changes
- Write simple commit messages prefixed using feat/chore/bug
- Before starting work, define the task and add it to TASKS.md
- Always add documentation for all cli arguments, config flags, etc.
- This tool will support many other cli tools with their different names and oddities for these types of files, we should include this in our code architecture such that there can be support for multiple "output formats"
- Always add unit tests for all functionality you add, aim for at least full branch coverage
- Write integration tests that test the full functionality, make sure XDG_CONFIG_HOME is overridden for the test execution such that the test-fragments can be part of the repository
- Always run linters with `--fix` or similar flags to automatically use safe fixes

## Git
- Always start new features in a feature branch `feat/<some-name>`
- Make many small commits

## Notion
- Only work in pages where the root page belongs to the project you are working in, usually named the same as the repository.

## NEVER DO
- NEVER EVER ADD CODE ATTRIBUTIONS IN COMMIT DESCRIPTIONS THAT REFERENCES THE CLI TOOL!

## Features
- Config location under XDG_CONFIG_HOME/.ctx/config.json
    - Config can be overridden using --config-file cli argument
- Files location for the fragments under XDG_CONFIG_HOME/.ctx/fragments
- Default tags to include should be configurable in the context file and possible to override via a suitable cli argument
- Config file should have a corresponding json schema describing it
- The main flow of the tool will be:
    - User runs `ctx build`
    - The tool visualizes available tags by grepping the tags from the fragments directory
    - The user can select tags (or these are preselected if --tags is provided)
    - Then the user can optionally deselect specific tags if needed
    - When the user is done the tool asks for confirmation and then combines the fragments
    - The tool also writes a "command file" for replication, i.e. a file that the tool could use to recreate the generated file
        - This file should include all fragments that was combined, this will allow the user to easily update their file if they update any of the fragment
- If the user does not have a default tool added in the config we need to ask the tool(s) they want to output for.
- The user should be able to select multiple output formats (for example opencode expects AGENTS.md and gemini expects GEMINI.md)
- The user should be able to select a custom output format as well, for example if they want a custom variant they can provide for specific agents
