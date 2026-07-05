package no_commented_out_tests_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_commented_out_tests"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoCommentedOutTestsRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&no_commented_out_tests.NoCommentedOutTestsRule,
		[]rule_tester.ValidTestCase{
			{Code: `// foo("bar", function () {})`},
			{Code: `describe("foo", function () {})`},
			{Code: `it("foo", function () {})`},
			{Code: `test("foo", function () {})`},
			{Code: `test.for([1])("foo", function () {})`},
			{Code: `describe.for([1])("foo", function () {})`},
			{Code: `test.runIf(flag)("foo", function () {})`},
			{Code: `describe.skipIf(flag)("foo", function () {})`},
			{Code: `// latest(dates)`},
			{Code: `// TODO: ship it soon`},
			{Code: `// fit("foo", function () {})`},
			{Code: `// xit("foo", function () {})`},
			{Code: `// xdescribe("foo", function () {})`},
			{Code: `// xtest("foo", function () {})`},
			{Code: `#!/usr/bin/env node`},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: `// describe("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: `// it.skip("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: `// test.only("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: `// test.for([1])("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: `// describe.runIf(flag)("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: `// test["skip"]("foo", function () {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: "// test(\n//   \"foo\", function () {}\n// )\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: "/* test\n  (\n    \"foo\", function () {}\n  )\n*/\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 1},
				},
			},
			{
				Code: "foo()\n/*\n  describe(\"has title but no callback\", () => {})\n*/\nbar()\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 2, Column: 1},
				},
			},
			{
				Code: "const 中文 = 1;\n// describe(\"x\", () => {})\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 2, Column: 1},
				},
			},
			{
				Code: "const e = \"🚀\";\n// test(\"x\", () => {})\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 2, Column: 1},
				},
			},
			{
				Code: `const 中文 = 1; // describe("a", () => {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 15},
				},
			},
			{
				Code: `const e = "🚀"; // test("x", () => {})`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "commentedTests", Line: 1, Column: 17},
				},
			},
		},
	)
}
