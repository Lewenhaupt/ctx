module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [2, 'always', [
      'feat',     // New feature
      'fix',      // Bug fix
      'docs',     // Documentation changes
      'style',    // Code style changes (formatting, etc.)
      'refactor', // Code refactoring
      'test',     // Adding or updating tests
      'chore'     // Maintenance tasks
    ]],
    'subject-case': [2, 'always', 'sentence-case'],
    'subject-empty': [2, 'never'],
    'subject-full-stop': [2, 'never', '.'],
    'subject-max-length': [2, 'always', 72],
    'body-leading-blank': [1, 'always'],
    'footer-leading-blank': [1, 'always'],
    'body-max-line-length': [0], // Disable body line length limit
    'footer-max-line-length': [0] // Disable footer line length limit
  },
  plugins: [
    {
      rules: {
        'no-unauthorized-attribution': (parsed) => {
          const message = parsed.raw.toLowerCase();
          
          // Check for code attributions
          const hasAttribution = message.includes('co-authored-by') || 
                               message.includes('authored-by') ||
                               message.includes('generated with');
          
          // Allow opencode attribution as it's required by the project
          const hasOpenCodeAttribution = message.includes('generated with [opencode]');
          
          if (hasAttribution && !hasOpenCodeAttribution) {
            return [false, 'Commit messages should not contain code attributions (except required opencode attribution)'];
          }
          return [true];
        }
      }
    }
  ]
};