package no_identical_title

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwNoIdenticalTitle "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_identical_title"
)

var NoIdenticalTitleRule = testfwNoIdenticalTitle.NewRule(testfwNoIdenticalTitle.Config{
	Name:                   "rstest/no-identical-title",
	Parse:                  rstestUtils.ParseRstestFnCall,
	ParameterizedModifiers: []string{"each", "for"},
})
