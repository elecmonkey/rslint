# expect-expect

## Rule Details

Ensure every Rstest test callback contains at least one assertion. Tests that
run code without verifying outcomes can pass while checking nothing.

The rule tracks Rstest `test` and `it` calls, including parameterized forms such
as `test.each` and `test.for`. `test.todo` is ignored because it intentionally
has no callback body.

Examples of **incorrect** code for this rule:

```ts
test('loads user', () => {
  loadUser();
});

it('does nothing', () => {});
```

Examples of **correct** code for this rule:

```ts
import { expect, test } from '@rstest/core';

test('loads user', () => {
  expect(loadUser()).toEqual({ id: 1 });
});

test('uses context expect', ({ expect }) => {
  expect(1 + 1).toBe(2);
});

test.todo('covers old behavior');
```

## Options

```ts
interface Options {
  assertFunctionNames?: string[];
  additionalTestBlockFunctions?: string[];
}
```

- `assertFunctionNames` (default `["expect"]`): callee chains that count as
  assertions. Wildcards follow the Jest rule behavior: `*` matches a segment and
  `**` matches zero or more segments.
- `additionalTestBlockFunctions`: extra function names treated like test block
  functions, useful for local wrappers.
