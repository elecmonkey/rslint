# no-mocks-import

## Rule Details

Disallow importing files from a `__mocks__` directory directly. Rstest resolves
manual mocks through `rs.mock()`; tests should import the original module path
and let Rstest replace it with the manual mock implementation.

Examples of **incorrect** code for this rule:

```ts
import api from '../src/__mocks__/api';
require('../src/__mocks__/api');
```

Examples of **correct** code for this rule:

```ts
import { rs } from '@rstest/core';
import api from '../src/api';

rs.mock('../src/api');
```

Directly importing `__mocks__` bypasses Rstest's module mocking mechanism and
can create confusing references that are different from the module mocked via
`rs.mock()`.
