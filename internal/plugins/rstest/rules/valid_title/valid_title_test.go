package valid_title_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/valid_title"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestValidTitleRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&valid_title.ValidTitleRule,
		[]rule_tester.ValidTestCase{
			{Code: `describe("suite", () => {});`},
			{Code: `test("works", () => {});`},
			{Code: `it("works", () => {});`},
			{Code: "test(`works`, () => {});"},
			{Code: `test("works " + suffix, () => {});`},
			{Code: `test.each([])("%s %d %i %f %j %o %O %c %# %$ %%", () => {});`},
			{Code: `test.for([])("%s %d %i %f %j %o %O %c %# %$ %%", () => {});`},
			{Code: `describe.each([])("%s %O", () => {});`},
			{Code: `describe.for([])("%c %$", () => {});`},
			{Code: `test.todo("todo title");`},
			{
				Code:    `test(123, () => {});`,
				Options: []interface{}{map[string]interface{}{"ignoreTypeOfTestName": true}},
			},
			{
				Code:    `describe(myFunction, () => {});`,
				Options: []interface{}{map[string]interface{}{"ignoreTypeOfDescribeName": true}},
			},
			{Code: `import { test as nodeTest } from "node:test"; nodeTest("", () => {});`},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: `describe("", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "emptyTitle", Line: 1, Column: 1},
				},
			},
			{
				Code: `test("", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "emptyTitle", Line: 1, Column: 1},
				},
			},
			{
				Code: `it(123, () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "titleMustBeString", Line: 1, Column: 4},
				},
			},
			{
				Code: `describe(myFunction, () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "titleMustBeString", Line: 1, Column: 10},
				},
			},
			{
				Code: `test(" test has surrounding spaces ", () => {});`,
				Output: []string{
					`test("test has surrounding spaces", () => {});`,
					`test("has surrounding spaces", () => {});`,
				},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "accidentalSpace", Line: 1, Column: 6},
				},
			},
			{
				Code:   `test("test works", () => {});`,
				Output: []string{`test("works", () => {});`},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "duplicatePrefix", Line: 1, Column: 6},
				},
			},
			{
				Code:   `it("it works", () => {});`,
				Output: []string{`it("works", () => {});`},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "duplicatePrefix", Line: 1, Column: 4},
				},
			},
			{
				Code:   `describe("describe works", () => {});`,
				Output: []string{`describe("works", () => {});`},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "duplicatePrefix", Line: 1, Column: 10},
				},
			},
			{
				Code: `test.each([])("%p", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "invalidEachSpecifier", Line: 1, Column: 15},
				},
			},
			{
				Code: `test.for([])("%p", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "invalidEachSpecifier", Line: 1, Column: 14},
				},
			},
			{
				Code:    `test("must have tag", () => {});`,
				Options: []interface{}{map[string]interface{}{"mustMatch": "#unit"}},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "mustMatch", Line: 1, Column: 6},
				},
			},
			{
				Code:    `test("has forbidden word", () => {});`,
				Options: []interface{}{map[string]interface{}{"disallowedWords": []interface{}{"forbidden"}}},
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "disallowedWord", Line: 1, Column: 6},
				},
			},
		},
	)
}
