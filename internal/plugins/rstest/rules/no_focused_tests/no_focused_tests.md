# no-focused-tests

## Rule Details

Disallow focused Rstest tests and suites. Focused tests such as `test.only` or
`describe.only` cause only part of a test file to run and can hide failures in
CI or local verification.

Examples of **incorrect** code for this rule:

```ts
describe.only('math', () => {});
test.only('adds numbers', () => {});
it.only('works', () => {});

test.concurrent.only('runs alone', () => {});
test.only.for([1, 2])('case %s', () => {});
describe.runIf(flag).only('suite', () => {});
```

Examples of **correct** code for this rule:

```ts
describe('math', () => {});
test('adds numbers', () => {});
it('works', () => {});

test.skip('temporarily disabled', () => {});
test.todo('covers old behavior');
test.concurrent('runs with peers', () => {});
test.for([1, 2])('case %s', () => {});
```

## Differences from Jest

Rstest does not expose Jest's focused aliases such as `fit` or `fdescribe`.
This rule only reports `.only` chains for Rstest. The Jest rule keeps reporting
both `.only` and Jest's `f*` aliases for Jest projects.
