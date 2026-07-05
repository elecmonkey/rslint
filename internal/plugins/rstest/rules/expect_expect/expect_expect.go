package expect_expect

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwExpectExpect "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/expect_expect"
)

var ExpectExpectRule = testfwExpectExpect.NewRule(testfwExpectExpect.Config{
	Name:  "rstest/expect-expect",
	Parse: rstestUtils.ParseRstestFnCall,
})
