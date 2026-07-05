package valid_title

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwValidTitle "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/valid_title"
)

var ValidTitleRule = testfwValidTitle.NewRule(testfwValidTitle.Config{
	Name:                   "rstest/valid-title",
	Parse:                  rstestUtils.ParseRstestFnCall,
	ParameterizedModifiers: []string{"each", "for"},
	AllowedEachSpecifiers:  "sdifjoOc#$%",
})
