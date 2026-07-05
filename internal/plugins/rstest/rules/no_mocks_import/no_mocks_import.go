package no_mocks_import

import testfwNoMocksImport "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_mocks_import"

var NoMocksImportRule = testfwNoMocksImport.NewRule(testfwNoMocksImport.Config{
	Name:    "rstest/no-mocks-import",
	MockAPI: "rs.mock",
})
