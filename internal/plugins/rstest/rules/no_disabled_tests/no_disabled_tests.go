package no_disabled_tests

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwNoDisabledTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_disabled_tests"
)

var NoDisabledTestsRule = testfwNoDisabledTests.NewRule(testfwNoDisabledTests.Config{
	Name:  "rstest/no-disabled-tests",
	Parse: rstestUtils.ParseRstestFnCall,
})
