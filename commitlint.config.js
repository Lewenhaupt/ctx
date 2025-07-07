module.exports = {
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
    'type-empty': [2, 'never'],
    'type-case': [2, 'always', 'lower-case'],
    'subject-empty': [2, 'never'],
    'subject-full-stop': [2, 'never', '.'],
    'subject-max-length': [2, 'always', 72],
    'header-max-length': [2, 'always', 72],
    'body-leading-blank': [1, 'always'],
    'footer-leading-blank': [1, 'always'],
    'body-max-line-length': [0], // Disable body line length limit
    'footer-max-line-length': [0], // Disable footer line length limit
    'scope-empty': [0, 'never']
  }
};
