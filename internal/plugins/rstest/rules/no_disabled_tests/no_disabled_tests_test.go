package no_disabled_tests_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_disabled_tests"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoDisabledTestsRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&no_disabled_tests.NoDisabledTestsRule,
		[]rule_tester.ValidTestCase{
			{Code: `describe("foo", function () {})`},
			{Code: `it("foo", function () {})`},
			{Code: `test("foo", function () {})`},
			{Code: `describe.only("foo", function () {})`},
			{Code: `it.only("foo", function () {})`},
			{Code: `test.only("foo", function () {})`},
			{Code: `test.todo("fill this later")`},
			{Code: `test("runtime skip", context => { context.skip(); })`},
			{Code: `test.for([1])("case %s", (value, context) => { context.expect(value).toBe(1); })`},
			{Code: `var appliedSkip = describe.skip; appliedSkip.apply(describe)`},
			{Code: `var calledSkip = it.skip; calledSkip.call(it)`},
			{Code: `fit("foo", function () {})`},
			{Code: `xit("foo", function () {})`},
			{
				Code: `
					import { test } from './test-utils';

					test('something');
				`,
			},
			{
				Code: `
					import { test as nodeTest } from 'node:test';

					nodeTest('something');
				`,
			},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: `describe.skip("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `describe.skip.each([1, 2, 3])("%s", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `it.skip("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `it["skip"]("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `test.skip("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `test.skip.for([1])("case %s", () => {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `test["skip"]("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "skippedTest", Line: 1, Column: 1},
				},
			},
			{
				Code: `it("has title but no callback")`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "missingFunction", Line: 1, Column: 1},
				},
			},
			{
				Code: `test("has title but no callback")`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "missingFunction", Line: 1, Column: 1},
				},
			},
			{
				Code: "import { test } from '@rstest/core';\n\ntest('something');",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "missingFunction", Line: 3, Column: 1},
				},
			},
		},
	)
}
