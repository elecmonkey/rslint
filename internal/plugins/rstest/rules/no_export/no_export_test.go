package no_export_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_export"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoExportRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&no_export.NoExportRule,
		[]rule_tester.ValidTestCase{
			{Code: `describe("a test", () => { expect(1).toBe(1); })`},
			{Code: `export const helper = "valid"`},
			{Code: `export default function () {}`},
			{Code: `module.exports = function(){}`},
			{Code: `module.exports.myThing = "valid";`},
			{Code: `module.export.foo = "valid"; test("a test", () => {});`},
			{Code: `const exports = "exports"; module[exports] = {}; test("a test", () => {});`},
			{Code: `const module = { exports: {} }; module.exports.foo = "valid"; test("a test", () => {});`},
			{Code: `const run = (module: { exports: object }) => { module.exports.foo = "valid" }; test("a test", () => {});`},
			{Code: `const exports = { foo: "" }; exports.foo = "valid"; test("a test", () => {});`},
			{Code: `const run = (exports: { foo: string }) => { exports.foo = "valid" }; test("a test", () => {});`},
			{Code: `import { test as nodeTest } from "node:test"; export const helper = "valid"; nodeTest("case", () => {});`},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: `export const helper = "invalid"; test("a test", () => { expect(1).toBe(1); });`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 33},
				},
			},
			{
				Code: `
export const helper = 'invalid';

test.for([1])('case %s', () => {
  expect(1).toBe(1);
});
`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 2, Column: 1, EndColumn: 33},
				},
			},
			{
				Code: `export default function() {}; test("a test", () => { expect(1).toBe(1); });`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 29},
				},
			},
			{
				Code: `export = function() {}; test("a test", () => { expect(1).toBe(1); });`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 24},
				},
			},
			{
				Code: `module.exports["invalid"] = function() {}; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 26},
				},
			},
			{
				Code: `module.exports = function() {}; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 15},
				},
			},
			{
				Code: `module["exports"] = function() {}; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 18},
				},
			},
			{
				Code: "module[`exports`].foo = function() {}; test(\"a test\", () => {});",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 22},
				},
			},
			{
				Code: `module.exports.foo.bar = function() {}; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 23},
				},
			},
			{
				Code: `module.exports ||= {}; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 15},
				},
			},
			{
				Code: `value = module.exports; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 9, EndColumn: 23},
				},
			},
			{
				Code: `exports.foo = "invalid"; test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 12},
				},
			},
			{
				Code: `export import foo = require("./foo"); test("a test", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 38},
				},
			},
			{
				Code: `export const helper = "invalid"; describe("a suite", () => {});`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "unexpectedExport", Line: 1, Column: 1, EndColumn: 33},
				},
			},
		},
	)
}
