package no_export

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwNoExport "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_export"
)

var NoExportRule = testfwNoExport.NewRule(testfwNoExport.Config{
	Name:  "rstest/no-export",
	Parse: rstestUtils.ParseRstestFnCall,
})
