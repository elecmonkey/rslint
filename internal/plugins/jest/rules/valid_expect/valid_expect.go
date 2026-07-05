package valid_expect

import (
	jestUtils "github.com/web-infra-dev/rslint/internal/plugins/jest/utils"
	testfwValidExpect "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/valid_expect"
)

var ValidExpectRule = testfwValidExpect.NewRule(testfwValidExpect.Config{
	Name:        "jest/valid-expect",
	ParseExpect: jestUtils.ParseJestExpectCallWithReason,
})
