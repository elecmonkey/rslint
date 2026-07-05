# no-identical-title

## Rule Details

Disallow duplicate static titles for tests or `describe` blocks in the same
scope. Duplicate titles make reporter output harder to read and make failures
harder to identify.

Only static string and simple template titles are checked. Dynamic titles and
parameterized Rstest calls such as `test.each`, `test.for`, `describe.each`,
and `describe.for` are ignored.

Examples of **incorrect** code for this rule:

```ts
test('loads user', () => {});
test('loads user', () => {});

describe('api', () => {});
describe('api', () => {});

describe('math', () => {
  it('adds numbers', () => {});
  it('adds numbers', () => {});
});
```

Examples of **correct** code for this rule:

```ts
test('loads user', () => {});
test('saves user', () => {});

describe('math', () => {
  it('adds numbers', () => {});

  describe('nested', () => {
    it('adds numbers', () => {});
  });
});

test.for([1, 2])('case %s', () => {});
test.for([3, 4])('case %s', () => {});
```
