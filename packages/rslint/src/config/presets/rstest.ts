import type { RslintConfigEntry } from '../define-config.js';

const recommended: RslintConfigEntry = {
  plugins: ['rstest'],
  rules: {
    'rstest/expect-expect': 'warn',
    'rstest/no-commented-out-tests': 'warn',
    'rstest/no-focused-tests': 'error',
    'rstest/no-identical-title': 'error',
    'rstest/no-mocks-import': 'error',
  },
};

export { recommended };
