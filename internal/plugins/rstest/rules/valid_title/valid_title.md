# valid-title

## Rule Details

Enforce valid titles on Rstest `describe`, `test`, and `it` blocks. Titles
should be readable string literals, should not be empty, and should not contain
accidental leading or trailing spaces.

For Rstest parameterized tests, this rule validates printf-style placeholders
in `test.each`, `test.for`, `describe.each`, and `describe.for` titles. Rstest
supports `%s`, `%d`, `%i`, `%f`, `%j`, `%o`, `%O`, `%c`, `%#`, `%$`, and `%%`.

Examples of **incorrect** code for this rule:

```ts
test('', () => {});
it(123, () => {});
describe(' describe users ', () => {});
test('test loads user', () => {});
test.for([1])('%p', () => {});
```

Examples of **correct** code for this rule:

```ts
test('loads user', () => {});
it('saves user', () => {});
describe('users', () => {});

test.for([1])('case %s', () => {});
test.each([{ name: 'Alice' }])('loads $name', () => {});
```

## Differences from Jest

Rstest uses `test.for` and `describe.for` in addition to `each`, so this rule
validates both forms. Rstest printf titles allow `%O` and `%c`, while Jest's
`%p` placeholder is not part of Rstest's `formatName()` behavior and is reported
as invalid.
