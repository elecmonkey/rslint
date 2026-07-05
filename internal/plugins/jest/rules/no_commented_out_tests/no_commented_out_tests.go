package no_commented_out_tests

import (
	"regexp"

	testfwNoCommentedOutTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_commented_out_tests"
)

// Port of eslint-plugin-jest no-commented-out-tests:
// /^\s*[xf]?(test|it|describe)(\.\w+|\[['"]\w+['"]\])?\s*\(/mu
var commentedTestRegexp = regexp.MustCompile(`(?m)^\s*[xf]?(test|it|describe)(\.\w+|\[['"]\w+['"]\])?\s*\(`)

var NoCommentedOutTestsRule = testfwNoCommentedOutTests.NewRule("jest/no-commented-out-tests", commentedTestRegexp)
