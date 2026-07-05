package no_disabled_tests

import (
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	"github.com/web-infra-dev/rslint/internal/rule"
)

type ParseFn func(node *ast.Node, ctx rule.RuleContext) *testfw.ParsedFnCall
type PendingCallFn func(node *ast.Node, ctx rule.RuleContext) bool

type Config struct {
	Name               string
	Parse              ParseFn
	PendingCall        PendingCallFn
	DisabledNamePrefix string
}

func buildErrorMissingFunctionMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "missingFunction",
		Description: "Test is missing function argument",
	}
}

func buildErrorSkippedTestMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "skippedTest",
		Description: "Tests should not be skipped",
	}
}

func NewRule(config Config) rule.Rule {
	return rule.Rule{
		Name: config.Name,
		Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
			return rule.RuleListeners{
				ast.KindCallExpression: func(node *ast.Node) {
					if config.PendingCall != nil && config.PendingCall(node, ctx) {
						ctx.ReportNode(node, buildErrorSkippedTestMessage())
						return
					}

					fnCall := config.Parse(node, ctx)
					if fnCall == nil ||
						fnCall.Kind != testfw.FnKindDescribe &&
							fnCall.Kind != testfw.FnKindTest {
						return
					}

					if (config.DisabledNamePrefix != "" && strings.HasPrefix(fnCall.Name, config.DisabledNamePrefix)) ||
						slices.Contains(fnCall.Members, "skip") {
						ctx.ReportNode(node, buildErrorSkippedTestMessage())
					}

					if fnCall.Kind == testfw.FnKindTest {
						if len(node.Arguments()) < 2 && !slices.Contains(fnCall.Members, "todo") {
							ctx.ReportNode(node, buildErrorMissingFunctionMessage())
						}
					}
				},
			}
		},
	}
}
