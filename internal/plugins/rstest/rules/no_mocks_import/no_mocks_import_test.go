package no_mocks_import_test

import (
	"testing"

	"github.com/web-infra-dev/rslint/internal/plugins/rstest/fixtures"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_mocks_import"
	"github.com/web-infra-dev/rslint/internal/rule_tester"
)

func TestNoMocksImportRule(t *testing.T) {
	rule_tester.RunRuleTester(
		fixtures.GetRootDir(),
		"tsconfig.json",
		t,
		&no_mocks_import.NoMocksImportRule,
		[]rule_tester.ValidTestCase{
			{Code: `import { rs } from "@rstest/core"; rs.mock("./module")`},
			{Code: `import thing from "thing"`},
			{Code: `require("somethingElse")`},
			{Code: `require("./__mocks__.js")`},
			{Code: `require("./__mocks__x")`},
			{Code: `require("./x__mocks__/x")`},
			{Code: `require()`},
			{Code: `var path = "./__mocks__.js"; require(path)`},
		},
		[]rule_tester.InvalidTestCase{
			{
				Code: `require("./__mocks__")`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noManualImport", Line: 1, Column: 9},
				},
			},
			{
				Code: `require("./__mocks__/index")`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noManualImport", Line: 1, Column: 9},
				},
			},
			{
				Code: `require("__mocks__/index")`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noManualImport", Line: 1, Column: 9},
				},
			},
			{
				Code: `import thing from "./__mocks__/index"`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noManualImport", Line: 1, Column: 1},
				},
			},
			{
				Code: `import { api } from "../src/__mocks__/api"`,
				Errors: []rule_tester.InvalidTestCaseError{
					{MessageId: "noManualImport", Line: 1, Column: 1},
				},
			},
		},
	)
}
