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

# Git commit-msg hook to validate commit messages according to commitlint rules
# This hook validates that commit messages follow the format: type: description
# where type is one of: feat, fix, docs, style, refactor, test, chore

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Read the commit message from the file passed as argument
COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Skip validation for merge commits, revert commits, and fixup commits
if echo "$COMMIT_MSG" | grep -qE "^(Merge|Revert|fixup!|squash!)"; then
    echo -e "${YELLOW}Skipping commit message validation for special commit type${NC}"
    exit 0
fi

# Function to validate commit message format
validate_commit_message() {
    local msg="$1"
    
    # Check if message is empty
    if [ -z "$msg" ]; then
        echo -e "${RED}❌ Commit message cannot be empty${NC}"
        return 1
    fi
    
    # Check format: type: description
    if ! echo "$msg" | grep -qE '^(feat|fix|docs|style|refactor|test|chore): .+'; then
        echo -e "${RED}❌ Invalid commit message format${NC}"
        echo -e "${RED}Current message: $msg${NC}"
        echo ""
        echo -e "${YELLOW}Expected format: 'type: description'${NC}"
        echo -e "${YELLOW}Allowed types: feat, fix, docs, style, refactor, test, chore${NC}"
        echo ""
        echo "Examples:"
        echo "  feat: add new user authentication system"
        echo "  fix: resolve memory leak in parser"
        echo "  docs: update README with installation instructions"
        echo "  chore: update dependencies"
        return 1
    fi
    
    # Extract the first line (header) for additional checks
    local header=$(echo "$msg" | head -n1)
    
    # Check header length (max 72 characters)
    if [ ${#header} -gt 72 ]; then
        echo -e "${RED}❌ Commit message header too long (${#header} characters, max 72)${NC}"
        echo -e "${RED}Header: $header${NC}"
        return 1
    fi
    
    # Check that header doesn't end with a period
    if echo "$header" | grep -q '\.$'; then
        echo -e "${RED}❌ Commit message header should not end with a period${NC}"
        echo -e "${RED}Header: $header${NC}"
        return 1
    fi
    
    # Check that type is lowercase
    local type=$(echo "$header" | cut -d':' -f1)
    if [ "$type" != "$(echo "$type" | tr '[:upper:]' '[:lower:]')" ]; then
        echo -e "${RED}❌ Commit type must be lowercase${NC}"
        echo -e "${RED}Found: $type${NC}"
        return 1
    fi
    
    # Check that there's a space after the colon
    if ! echo "$header" | grep -qE '^[a-z]+: '; then
        echo -e "${RED}❌ Missing space after colon in commit message${NC}"
        echo -e "${RED}Header: $header${NC}"
        echo -e "${YELLOW}Expected format: 'type: description' (note the space after colon)${NC}"
        return 1
    fi
    
    return 0
}

# Validate the commit message
if validate_commit_message "$COMMIT_MSG"; then
    echo -e "${GREEN}✅ Commit message format is valid${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}Commit rejected due to invalid message format.${NC}"
    echo -e "${YELLOW}Please fix your commit message and try again.${NC}"
    exit 1
fi
EOF

# Make the hook executable
chmod +x "$HOOK_FILE"

echo -e "${GREEN}✅ Git commit-msg hook has been set up successfully!${NC}"
echo ""
echo "The hook will now validate commit messages according to the commitlint rules:"
echo "- Format: 'type: description'"
echo "- Allowed types: feat, fix, docs, style, refactor, test, chore"
echo "- Max header length: 72 characters"
echo "- No period at the end of the header"
echo ""
echo "Example valid commit messages:"
echo "  feat: add new user authentication system"
echo "  fix: resolve memory leak in parser"
echo "  docs: update README with installation instructions"
echo "  chore: update dependencies"