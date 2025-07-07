#!/usr/bin/env bash

# Script to set up git hooks for the ctx project
# This script should be run by developers after cloning the repository

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Setting up git hooks for ctx project...${NC}"

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "Error: This script must be run from the root of the git repository"
    exit 1
fi

# Remove Husky hooks path if it exists
if git config --get core.hooksPath >/dev/null 2>&1; then
    echo "Removing existing hooks path configuration..."
    git config --unset core.hooksPath
fi

# Create the commit-msg hook
HOOK_FILE=".git/hooks/commit-msg"

cat > "$HOOK_FILE" << 'EOF'
#!/usr/bin/env bash

# Git commit-msg hook using commitlint CLI
# Uses the same commitlint configuration as the GitHub workflow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

COMMIT_MSG_FILE="$1"

# Skip validation for merge commits, revert commits, and fixup commits
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")
if echo "$COMMIT_MSG" | grep -qE "^(Merge|Revert|fixup!|squash!)"; then
    echo -e "${YELLOW}Skipping commit message validation for special commit type${NC}"
    exit 0
fi

# Find commitlint binary
COMMITLINT_CMD=""

# Check if commitlint is available globally
if command -v commitlint >/dev/null 2>&1; then
    COMMITLINT_CMD="commitlint"
# Check if commitlint is available via npx
elif command -v npx >/dev/null 2>&1; then
    COMMITLINT_CMD="npx @commitlint/cli"
# Try to install commitlint globally if npm is available
elif command -v npm >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  commitlint not found. Installing globally...${NC}"
    if npm install -g @commitlint/cli >/dev/null 2>&1; then
        COMMITLINT_CMD="commitlint"
    else
        echo -e "${YELLOW}⚠️  Global install failed. Using npx...${NC}"
        COMMITLINT_CMD="npx @commitlint/cli"
    fi
else
    echo -e "${RED}❌ Neither commitlint nor npm/npx found.${NC}"
    echo "Please install Node.js and npm, then install commitlint:"
    echo "npm install -g @commitlint/cli"
    echo ""
    echo "Or ensure commitlint is available in your PATH."
    exit 1
fi

# Run commitlint on the commit message
if $COMMITLINT_CMD --edit "$COMMIT_MSG_FILE" --verbose; then
    echo -e "${GREEN}✅ Commit message format is valid${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Commit message validation failed${NC}"
    echo -e "${YELLOW}Please fix your commit message according to the commitlint rules and try again.${NC}"
    echo ""
    echo "Expected format: 'type: description'"
    echo "Allowed types: feat, fix, docs, style, refactor, test, chore"
    echo ""
    echo "Examples:"
    echo "  feat: add new user authentication system"
    echo "  fix: resolve memory leak in parser"
    echo "  docs: update README with installation instructions"
    echo "  chore: update dependencies"
    exit 1
fi
EOF

# Make the hook executable
chmod +x "$HOOK_FILE"

echo -e "${GREEN}✅ Git commit-msg hook has been set up successfully!${NC}"
echo ""
echo "The hook uses commitlint CLI to validate commit messages using the same"
echo "configuration as the GitHub workflow (commitlint.config.js)."
echo ""
echo "Commit message format: 'type: description'"
echo "Allowed types: feat, fix, docs, style, refactor, test, chore"
echo ""
echo "Example valid commit messages:"
echo "  feat: add new user authentication system"
echo "  fix: resolve memory leak in parser"
echo "  docs: update README with installation instructions"
echo "  chore: update dependencies"
echo ""
echo "Note: commitlint will be installed globally if not already available."