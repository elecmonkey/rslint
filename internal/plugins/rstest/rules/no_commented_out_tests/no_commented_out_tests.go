package no_commented_out_tests

import (
	"regexp"

	testfwNoCommentedOutTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_commented_out_tests"
)

// Rstest test APIs are test, it and describe. Unlike Jest, Rstest does not
// expose x*/f* aliases, so this intentionally only matches the canonical names.
var commentedTestRegexp = regexp.MustCompile(`(?m)^\s*(test|it|describe)(\.\w+|\[['"]\w+['"]\])?\s*\(`)

var NoCommentedOutTestsRule = testfwNoCommentedOutTests.NewRule("rstest/no-commented-out-tests", commentedTestRegexp)
