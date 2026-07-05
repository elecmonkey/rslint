package no_focused_tests

import (
	jestUtils "github.com/web-infra-dev/rslint/internal/plugins/jest/utils"
	testfwNoFocusedTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_focused_tests"
)

var NoFocusedTestsRule = testfwNoFocusedTests.NewRule(testfwNoFocusedTests.Config{
	Name:              "jest/no-focused-tests",
	Parse:             jestUtils.ParseJestTestFnCall,
	FocusedNamePrefix: "f",
	FocusedReplacement: map[string]string{
		"fdescribe": "describe",
		"fit":       "it",
	},
})
