# GitHub Actions Workflows

This directory contains the GitHub Actions workflows for the ctx project.

## Workflows

### 🔄 CI (`ci.yml`)
**Triggers**: Push to main, Pull requests
- **Go Testing**: Runs tests on multiple Go versions (1.21, 1.22, 1.23)
- **Linting**: golangci-lint, gofmt, go vet
- **Nix Build**: Tests Nix package builds
- **Cross-platform Build**: Tests compilation for multiple platforms

### 📝 Commit Linting (`lint-commits.yml`)
**Triggers**: Push to main, Pull requests
- Validates commit message format: `type: description`
- Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Enforces project-specific rules (no unauthorized code attributions)

### 🚀 Release (`release.yml`)
**Triggers**: Tag push (`v*.*.*`)
- **Build Binaries**: Cross-compiles for Linux, macOS, Windows (amd64/arm64)
- **Create Release**: Auto-generates changelog and uploads binaries
- **Update Flake**: Creates PR to update flake.nix with new version and vendorHash

### 📦 Dependencies (`deps.yml`)
**Triggers**: Weekly schedule (Sundays), Manual dispatch
- **Go Dependencies**: Updates Go modules and creates PR
- **Nix Dependencies**: Updates flake inputs and creates PR
- **Security Audit**: Runs vulnerability scans

## Configuration Files

### `.golangci.yml`
Comprehensive Go linting configuration with:
- Code quality checks
- Security scanning
- Style enforcement
- Performance optimizations

### `commitlint.config.js`
Commit message validation rules:
- Enforces conventional commit format
- Project-specific attribution rules
- Message length limits

### `dependabot.yml`
Automated dependency updates for:
- Go modules (weekly)
- GitHub Actions (weekly)

## Release Process

1. **Create and push a tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Automated workflow**:
   - Builds binaries for all platforms
   - Creates GitHub release with changelog
   - Uploads release assets
   - Creates PR to update flake.nix

3. **Manual steps**:
   - Review and merge the flake.nix update PR
   - Announce the release

## Security

- All workflows use pinned action versions
- Minimal permissions (GITHUB_TOKEN only)
- Vulnerability scanning included
- Dependency review on PRs

## Caching

- Go modules cache for faster builds
- Nix store cache (requires CACHIX_AUTH_TOKEN secret)
- Cross-platform build artifacts

## Required Secrets

- `GITHUB_TOKEN`: Automatically provided by GitHub
- `CACHIX_AUTH_TOKEN`: Optional, for Nix cache (improves build speed)

## Customization

To customize for your repository:

1. Update `dependabot.yml` with your GitHub username
2. Add your username to workflow reviewers/assignees
3. Configure Cachix if desired for faster Nix builds
4. Adjust Go version matrix in CI if needed