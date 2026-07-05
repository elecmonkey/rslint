package no_focused_tests

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwNoFocusedTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_focused_tests"
)

var NoFocusedTestsRule = testfwNoFocusedTests.NewRule(testfwNoFocusedTests.Config{
	Name:  "rstest/no-focused-tests",
	Parse: rstestUtils.ParseRstestFnCall,
})
