package no_identical_title_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_identical_title"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoIdenticalTitleRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&no_identical_title.NoIdenticalTitleRule,
		[]rule_tester.ValidTestCase{
			{Code: "it(); it();"},
			{Code: "describe(); describe();"},
			{Code: "describe('foo', () => {}); it('foo', () => {});"},
			{Code: "it('one', () => {});\nit('two', () => {});\n"},
			{Code: "describe('foo', () => {});\ndescribe('foe', () => {});\n"},
			{Code: "describe('foo', () => {\n  it('works', () => {});\n  describe('nested', () => {\n    it('works', () => {});\n  });\n});\n"},
			{Code: "test('number' + n, () => {});\ntest('number' + n, () => {});\n"},
			{Code: "it(`${n}`, () => {});\nit(`${n}`, () => {});\n"},
			{Code: "test.each([])('same title', () => {});\ntest.each([])('same title', () => {});\n"},
			{Code: "test.for([])('same title', () => {});\ntest.for([])('same title', () => {});\n"},
			{Code: "describe.each([])('same title', () => {});\ndescribe.each([])('same title', () => {});\n"},
			{Code: "describe.for([])('same title', () => {});\ndescribe.for([])('same title', () => {});\n"},
			{Code: "const test = { content: () => 'foo' };\ntest.content('same', () => {});\ntest.content('same', () => {});\n"},
			{Code: "import { test as nodeTest } from 'node:test';\nnodeTest('same', () => {});\nnodeTest('same', () => {});\n"},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: "it('works', () => {});\nit('works', () => {});\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "multipleTestTitle", Line: 2, Column: 4},
				},
			},
			{
				Code: "test('works', () => {});\ntest('works', () => {});\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "multipleTestTitle", Line: 2, Column: 6},
				},
			},
			{
				Code: "test.only('works', () => {});\ntest.concurrent('works', () => {});\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "multipleTestTitle", Line: 2, Column: 17},
				},
			},
			{
				Code: "describe('foo', () => {});\ndescribe('foo', () => {});\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "multipleDescribeTitle", Line: 2, Column: 10},
				},
			},
			{
				Code: "describe('foo', () => {\n  test('same', () => {});\n  test('same', () => {});\n});\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "multipleTestTitle", Line: 3, Column: 8},
				},
			},
			{
				Code: "import { test as rstestTest } from '@rstest/core';\nrstestTest('same', () => {});\nrstestTest('same', () => {});\n",
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "multipleTestTitle", Line: 3, Column: 12},
				},
			},
		},
	)
}
