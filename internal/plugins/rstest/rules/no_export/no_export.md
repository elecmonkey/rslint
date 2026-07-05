# no-export

## Rule Details

Disallow exporting from files that contain Rstest tests or suites. A test file
should be executed by Rstest as a test entry, not imported by another module as a
shared API.

This rule reports exports only when the same file also contains a Rstest
`test`, `it`, or `describe` block.

It checks:

- ES module exports such as `export const`, `export default`, and `export =`
- CommonJS assignments rooted at `module.exports` or `exports`

Examples of **incorrect** code for this rule:

```ts
export const helper = 'shared';

test('uses helper', () => {
  expect(helper).toBe('shared');
});
```

```ts
module.exports = {};

describe('suite', () => {});
```

Examples of **correct** code for this rule:

```ts
const helper = 'local';

test('uses helper', () => {
  expect(helper).toBe('local');
});
```

```ts
export const helper = 'shared';
```

If helpers must be shared, move them to a dedicated helper module that does not
define tests.
