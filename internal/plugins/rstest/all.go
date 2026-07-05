package rstest

import (
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/expect_expect"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_commented_out_tests"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_disabled_tests"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_export"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_focused_tests"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_identical_title"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/no_mocks_import"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/valid_expect"
	"github.com/web-infra-dev/rslint/internal/plugins/rstest/rules/valid_title"
	"github.com/web-infra-dev/rslint/internal/rule"
)

func GetAllRules() []rule.Rule {
	return []rule.Rule{
		expect_expect.ExpectExpectRule,
		no_commented_out_tests.NoCommentedOutTestsRule,
		no_disabled_tests.NoDisabledTestsRule,
		no_export.NoExportRule,
		no_focused_tests.NoFocusedTestsRule,
		no_identical_title.NoIdenticalTitleRule,
		no_mocks_import.NoMocksImportRule,
		valid_expect.ValidExpectRule,
		valid_title.ValidTitleRule,
	}
}
