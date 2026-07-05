package expect_expect_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/expect_expect"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestExpectExpectRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&expect_expect.ExpectExpectRule,
		[]rule_tester.ValidTestCase{
			{Code: `test.todo("will test something eventually")`},
			{Code: `test("passes", () => expect(true).toBeDefined())`},
			{Code: `it("passes", () => expect(true).toBeDefined())`},
			{Code: `test("passes with context expect", ({ expect }) => expect(true).toBeDefined())`},
			{Code: `test.for([1])("case %s", (value, { expect }) => expect(value).toBe(1))`},
			{Code: `test.each([1])("case %s", value => expect(value).toBe(1))`},
			{Code: `test("passes", helper); function helper() { expect(true).toBeDefined() }`},
			{
				Code: `
        import { test as rstestTest, expect } from '@rstest/core';

        rstestTest('passes', () => {
          expect(true).toBeDefined();
        });
      `,
			},
			{
				Code: `
        test('uses custom assertion helper', () => {
          assertOk();
        });
      `,
				Options: []interface{}{map[string]interface{}{"assertFunctionNames": []interface{}{"assertOk"}}},
			},
			{
				Code: `
        scenario('custom block', () => {
          expect(true).toBeDefined();
        });
      `,
				Options: []interface{}{map[string]interface{}{"additionalTestBlockFunctions": []interface{}{"scenario"}}},
			},
			{
				Code: `
        import { test as nodeTest } from 'node:test';

        nodeTest('not rstest', () => {});
      `,
			},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: `test("fails", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noAssertions", Line: 1, Column: 1, EndLine: 1, EndColumn: 5},
				},
			},
			{
				Code: `it("fails", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noAssertions", Line: 1, Column: 1, EndLine: 1, EndColumn: 3},
				},
			},
			{
				Code: `test.skip("fails", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noAssertions", Line: 1, Column: 1, EndLine: 1, EndColumn: 10},
				},
			},
			{
				Code: `test.for([1])("case %s", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noAssertions", Line: 1, Column: 1, EndLine: 1, EndColumn: 14},
				},
			},
			{
				Code: `
        import { test as rstestTest } from '@rstest/core';

        rstestTest('fails', () => {});
      `,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noAssertions", Line: 4, Column: 9, EndLine: 4, EndColumn: 19},
				},
			},
			{
				Code: `
        scenario('custom block', () => {});
      `,
				Options: []interface{}{map[string]interface{}{"additionalTestBlockFunctions": []interface{}{"scenario"}}},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noAssertions", Line: 2, Column: 9, EndLine: 2, EndColumn: 17},
				},
			},
		},
	)
}
