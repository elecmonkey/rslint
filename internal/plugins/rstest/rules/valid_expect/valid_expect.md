# valid-expect

## Rule Details

Enforce valid Rstest `expect` usage. This rule checks that assertions have the
expected number of arguments, end with a called matcher, use valid modifiers,
and await or return async assertions.

Examples of **incorrect** code for this rule:

```ts
expect();
expect(value);
expect(value).toBeDefined;
expect(value).unknown.toBeDefined();

expect(Promise.resolve(1)).resolves.toBe(1);
expect.soft(value);
expect.poll(() => value).toBe(1);
```

Examples of **correct** code for this rule:

```ts
expect(value).toBeDefined();
expect.soft(value).toBeDefined();

await expect(Promise.resolve(1)).resolves.toBe(1);
return expect(Promise.resolve(1)).resolves.toBe(1);

await expect.poll(() => getStatus()).toBe('ready');
```

## Rstest-specific behavior

`expect.soft(value)` is treated like `expect(value)`: it must have valid
arguments and a called matcher.

`expect.poll(fn)` is an async assertion factory. Its matcher chain must be
awaited or returned.
