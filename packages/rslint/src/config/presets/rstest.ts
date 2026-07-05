import type { RslintConfigEntry } from '../define-config.js';

const recommended: RslintConfigEntry = {
  plugins: ['rstest'],
  rules: {
    'rstest/expect-expect': 'warn',
    'rstest/no-commented-out-tests': 'warn',
    'rstest/no-disabled-tests': 'warn',
    'rstest/no-export': 'error',
    'rstest/no-focused-tests': 'error',
    'rstest/no-identical-title': 'error',
    'rstest/no-mocks-import': 'error',
    'rstest/valid-expect': 'error',
    'rstest/valid-title': 'error',
  },
};

export { recommended };
