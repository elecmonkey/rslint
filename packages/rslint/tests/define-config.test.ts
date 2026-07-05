import { describe, test, expect } from '@rstest/core';
import { globalIgnores } from '../src/config/define-config.js';
import { rstestPlugin } from '../src/index.js';

describe('globalIgnores', () => {
  test('returns a config entry containing only the ignores', () => {
    expect(globalIgnores(['dist/**', 'node_modules/**'])).toEqual({
      ignores: ['dist/**', 'node_modules/**'],
    });
  });

  test('produces an entry with no keys other than `ignores`', () => {
    // The global-ignore semantics rely on the entry having ONLY `ignores`
    // (no files/rules/plugins/settings/languageOptions). Lock that in so the
    // go-side `isGlobalIgnoreEntry` keeps recognizing it as a global ignore.
    expect(Object.keys(globalIgnores(['dist/**']))).toEqual(['ignores']);
  });

  test('returns the same patterns array reference', () => {
    const patterns = ['dist/**'];
    expect(globalIgnores(patterns).ignores).toBe(patterns);
  });

  test('throws TypeError when patterns is not an array', () => {
    // @ts-expect-error testing the runtime guard against non-array input
    expect(() => globalIgnores('dist/**')).toThrow(TypeError);
    // @ts-expect-error testing the runtime guard against non-array input
    expect(() => globalIgnores('dist/**')).toThrow(
      'ignorePatterns must be an array',
    );
  });

  test('throws TypeError when patterns is empty', () => {
    expect(() => globalIgnores([])).toThrow(TypeError);
    expect(() => globalIgnores([])).toThrow(
      'ignorePatterns must contain at least one pattern',
    );
  });
});

describe('rstestPlugin preset', () => {
  test('exports the Rstest recommended preset', () => {
    expect(rstestPlugin.configs.recommended).toEqual({
      plugins: ['rstest'],
      rules: {
        'rstest/expect-expect': 'warn',
        'rstest/no-commented-out-tests': 'warn',
        'rstest/no-disabled-tests': 'warn',
        'rstest/no-focused-tests': 'error',
        'rstest/no-identical-title': 'error',
        'rstest/no-mocks-import': 'error',
      },
    });
  });
});
