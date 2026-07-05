package no_focused_tests_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_focused_tests"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoFocusedTestsRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&no_focused_tests.NoFocusedTestsRule,
		[]rule_tester.ValidTestCase{
			{Code: "describe()"},
			{Code: "it()"},
			{Code: "test()"},
			{Code: "describe.skip()"},
			{Code: "it.skip()"},
			{Code: "test.skip()"},
			{Code: "test.todo()"},
			{Code: "test.each()()"},
			{Code: "test.for()()"},
			{Code: "describe.each()()"},
			{Code: "describe.for()()"},
			{Code: "test.concurrent()"},
			{Code: "test.sequential()"},
			{Code: "test.runIf(flag)()"},
			{Code: "test.skipIf(flag)()"},
			{Code: "var appliedOnly = describe.only; appliedOnly.apply(describe)"},
			{Code: "var calledOnly = it.only; calledOnly.call(it)"},
			{Code: "fit()"},
			{Code: "fdescribe()"},
			{Code: "import { test as focusedTest } from 'node:test';\nfocusedTest.only()"},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: "describe.only()",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    10,
						EndLine:   1,
						EndColumn: 14,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "describe()",
							},
						},
					},
				},
			},
			{
				Code: "it.only()",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    4,
						EndLine:   1,
						EndColumn: 8,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "it()",
							},
						},
					},
				},
			},
			{
				Code: "test.only()",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    6,
						EndLine:   1,
						EndColumn: 10,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "test()",
							},
						},
					},
				},
			},
			{
				Code: "test.concurrent.only.each()()",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    17,
						EndLine:   1,
						EndColumn: 21,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "test.concurrent.each()()",
							},
						},
					},
				},
			},
			{
				Code: "test.only.for([1])('case', () => {})",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    6,
						EndLine:   1,
						EndColumn: 10,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "test.for([1])('case', () => {})",
							},
						},
					},
				},
			},
			{
				Code: "test.runIf(flag).only()",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    18,
						EndLine:   1,
						EndColumn: 22,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "test.runIf(flag)()",
							},
						},
					},
				},
			},
			{
				Code: `describe["only"]()`,
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      1,
						Column:    10,
						EndLine:   1,
						EndColumn: 16,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "describe()",
							},
						},
					},
				},
			},
			{
				Code: "import { test as rstestTest } from '@rstest/core';\nrstestTest.only()",
				Errors: []rule_tester.InvalidTestCaseError{
					{
						MessageId: "focusedTest",
						Line:      2,
						Column:    12,
						EndLine:   2,
						EndColumn: 16,
						Suggestions: []rule_tester.InvalidTestCaseSuggestion{
							{
								MessageId: "suggestRemoveFocus",
								Output:    "import { test as rstestTest } from '@rstest/core';\nrstestTest()",
							},
						},
					},
				},
			},
		},
	)
}
