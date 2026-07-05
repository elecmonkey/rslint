# no-disabled-tests

## Rule Details

Disallow skipped or incomplete Rstest tests. This rule reports tests and suites
marked with `.skip`, and reports `test` / `it` calls that omit the callback
function. Explicit `test.todo(...)` calls are allowed because they are visible in
test reports.

Examples of **incorrect** code for this rule:

```ts
describe.skip('suite', () => {});
test.skip('case', () => {});
it.skip('case', () => {});

test('missing callback');
it('missing callback');
```

Examples of **correct** code for this rule:

```ts
describe('suite', () => {});
test('case', () => {});
it('case', () => {});

test.todo('covers old behavior');

test('skip at runtime when needed', (context) => {
  context.skip();
});
```

## Differences from Jest

Rstest does not expose Jest's `xdescribe`, `xit`, or `xtest` aliases and does
not use Jasmine's `pending()` API. This rule therefore only reports Rstest's
`.skip` chains and missing test callbacks.
