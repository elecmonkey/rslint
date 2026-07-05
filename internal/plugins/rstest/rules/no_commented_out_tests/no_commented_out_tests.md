# no-commented-out-tests

## Rule Details

Disallow commenting out Rstest tests. Commented-out tests do not appear in test
reports and are easy to leave behind during review. Prefer deleting obsolete
tests, extracting reusable helpers, or using explicit Rstest APIs such as
`test.skip` or `test.todo`.

rslint walks each comment body line by line. If a line looks like a commented
Rstest `test`, `it`, or `describe` call, including chained forms such as
`.skip`, `.only`, `.each`, `.for`, `.runIf`, or bracket access such as
`test["skip"]`, the rule reports the whole comment range.

Examples of **incorrect** code for this rule:

```ts
// describe('math', () => {});
// it('adds numbers', () => {});
// test('loads user', () => {});

// test.skip('temporary disabled', () => {});
// test.for([1, 2])('case %s', () => {});

/*
test('old behavior', () => {
  expect(run()).toBe('ok');
});
*/
```

Examples of **correct** code for this rule:

```ts
import { describe, expect, test } from '@rstest/core';

describe('math', () => {
  test('adds numbers', () => {
    expect(1 + 1).toBe(2);
  });
});

test.todo('covers old behavior');

// latest(dates)
```

## Differences from Jest

Rstest does not expose Jest's `fit`, `fdescribe`, `xit`, `xtest`, or
`xdescribe` aliases, so this rule intentionally does not report those commented
forms. The Jest rule keeps reporting them for Jest projects.
