package rstest

import (
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/expect_expect"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_commented_out_tests"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_focused_tests"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_mocks_import"
	"github.com/web-infra-dev/rslint/internal/rule"
)

func GetAllRules() []rule.Rule {
	return []rule.Rule{
		expect_expect.ExpectExpectRule,
		no_commented_out_tests.NoCommentedOutTestsRule,
		no_focused_tests.NoFocusedTestsRule,
		no_mocks_import.NoMocksImportRule,
	}
}
