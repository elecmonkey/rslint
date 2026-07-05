package no_focused_tests

import (
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/microsoft/typescript-go/shim/core"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	"github.com/web-infra-dev/rslint/internal/rule"
	rslintUtils "github.com/web-infra-dev/rslint/internal/utils"
)

type ParseFn func(node *ast.Node, ctx rule.RuleContext) *testfw.ParsedFnCall

type Config struct {
	Name               string
	Parse              ParseFn
	FocusedNamePrefix  string
	FocusedReplacement map[string]string
}

func buildErrorFocusedTestMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "focusedTest",
		Description: "Unexpected focused test",
	}
}

func buildErrorSuggestRemoveFocusMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "suggestRemoveFocus",
		Description: "Suggest removing focus from test",
	}
}

func NewRule(config Config) rule.Rule {
	return rule.Rule{
		Name: config.Name,
		Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
			return rule.RuleListeners{
				ast.KindCallExpression: func(node *ast.Node) {
					fnCall := config.Parse(node, ctx)
					if fnCall == nil ||
						(fnCall.Kind != testfw.FnKindDescribe &&
							fnCall.Kind != testfw.FnKindTest) {
						return
					}

					if config.FocusedNamePrefix != "" && strings.HasPrefix(fnCall.Name, config.FocusedNamePrefix) {
						reportFocusedPrefixCall(ctx, node, fnCall, config)
						return
					}

					idx := slices.IndexFunc(fnCall.MemberEntries, func(entry testfw.MemberEntry) bool {
						return entry.Name == "only"
					})
					if idx < 0 {
						return
					}

					entry := fnCall.MemberEntries[idx]
					startRange := entry.Node.Loc.Pos() - 1
					endRange := entry.Node.Loc.End()
					if entry.Node.Kind != ast.KindIdentifier {
						endRange = entry.Node.End() + 1
					}

					ctx.ReportNodeWithSuggestions(
						entry.Node,
						buildErrorFocusedTestMessage(),
						rule.RuleSuggestion{
							Message: buildErrorSuggestRemoveFocusMessage(),
							FixesArr: []rule.RuleFix{
								rule.RuleFixRemoveRange(core.NewTextRange(startRange, endRange)),
							},
						},
					)
				},
			}
		},
	}
}

func reportFocusedPrefixCall(ctx rule.RuleContext, node *ast.Node, fnCall *testfw.ParsedFnCall, config Config) {
	callExpr := node.AsCallExpression()
	if callExpr == nil {
		return
	}

	callee := ast.SkipParentheses(callExpr.Expression)
	if callee == nil {
		return
	}

	reportNode := fnCall.Head.Local.Node
	if reportNode == nil {
		reportNode = callee
	}

	if fnCall.Head.Type == testfw.ImportedMode && fnCall.Name != fnCall.Head.Local.Value {
		ctx.ReportNode(reportNode, buildErrorFocusedTestMessage())
		return
	}

	replacement, ok := config.FocusedReplacement[fnCall.Name]
	if !ok {
		ctx.ReportNode(reportNode, buildErrorFocusedTestMessage())
		return
	}

	reportRange := rslintUtils.TrimNodeTextRange(ctx.SourceFile, reportNode)
	ctx.ReportNodeWithSuggestions(
		reportNode,
		buildErrorFocusedTestMessage(),
		rule.RuleSuggestion{
			Message: buildErrorSuggestRemoveFocusMessage(),
			FixesArr: []rule.RuleFix{
				rule.RuleFixReplaceRange(reportRange, replacement),
			},
		},
	)
}
