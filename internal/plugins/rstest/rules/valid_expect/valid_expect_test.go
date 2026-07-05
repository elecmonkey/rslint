package valid_expect_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/valid_expect"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestValidExpectRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&valid_expect.ValidExpectRule,
		[]rule_tester.ValidTestCase{
			{Code: "expect.hasAssertions()"},
			{Code: "expect.assertions(1)"},
			{Code: "expect.extend({})"},
			{Code: "expect.any(String)"},
			{Code: "expect.objectContaining({})"},
			{Code: "expect.unreachable()"},
			{Code: "expect('something').toEqual('something');"},
			{Code: "expect(true).not.toBe(false);"},
			{Code: "expect.soft(true).toBe(true);"},
			{Code: "test('valid-expect', async () => { await expect(Promise.resolve(2)).resolves.toBe(2); });"},
			{Code: "test('valid-expect', () => { return expect(Promise.resolve(2)).resolves.toBe(2); });"},
			{Code: "test('valid-expect', async () => { await expect.poll(() => 1).toBe(1); });"},
			{Code: "test('valid-expect', () => expect(Promise.resolve(2)).resolves.toBe(2));"},
			{Code: "import { expect as rstestExpect } from '@rstest/core'; rstestExpect(true).toBe(true);"},
			{Code: "import { expect } from 'chai'; expect(foo).to.equal(bar);"},
			{Code: "expect().pass();", Options: map[string]interface{}{"minArgs": 0}},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: "expect().toBe(2);",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "notEnoughArgs"},
				},
			},
			{
				Code: "expect('something');",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "matcherNotFound", Column: 1},
				},
			},
			{
				Code: "expect(true).toBeDefined;",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "matcherNotCalled", Column: 14},
				},
			},
			{
				Code: "expect(true).nope.toBeDefined();",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "modifierUnknown", Column: 1},
				},
			},
			{
				Code: "expect(true).not.resolves.toBeDefined();",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "modifierUnknown", Column: 1},
				},
			},
			{
				Code: "expect('something', 'else').toEqual('something');",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "tooManyArgs"},
				},
			},
			{
				Code: "expect.soft();",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "matcherNotFound"},
				},
			},
			{
				Code: "expect.soft(true);",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "matcherNotFound"},
				},
			},
			{
				Code: "expect(Promise.resolve(2)).resolves.toBe(2);",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "asyncMustBeAwaited", Column: 1},
				},
			},
			{
				Code:   "test('valid-expect', async () => { expect.poll(() => 1).toBe(1); });",
				Output: []string{"test('valid-expect', async () => { await expect.poll(() => 1).toBe(1); });"},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "asyncMustBeAwaited"},
				},
			},
			{
				Code: "test('valid-expect', async () => { await expect.poll(() => 1).resolves.toBe(1); });",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unsupportedPollMember"},
				},
			},
			{
				Code: "test('valid-expect', async () => { await expect.poll(() => () => {}).toThrow(); });",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unsupportedPollMember"},
				},
			},
			{
				Code: "test('valid-expect', async () => { await expect.poll(() => 'x').toMatchSnapshot(); });",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unsupportedPollMember"},
				},
			},
		},
	)
}
