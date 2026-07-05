package no_identical_title

import (
	"slices"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	"github.com/web-infra-dev/rslint/internal/rule"
)

type ParseFn func(node *ast.Node, ctx rule.RuleContext) *testfw.ParsedFnCall

type Config struct {
	Name                   string
	Parse                  ParseFn
	ParameterizedModifiers []string
}

func buildMultipleTestTitleMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "multipleTestTitle",
		Description: "Test title is used multiple times in the same describe block",
	}
}

func buildMultipleDescribeTitleMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "multipleDescribeTitle",
		Description: "Describe block title is used multiple times in the same describe block",
	}
}

type titleLayer struct {
	describeTitles map[string]struct{}
	testTitles     map[string]struct{}
}

func newTitleLayer() *titleLayer {
	return &titleLayer{
		describeTitles: make(map[string]struct{}),
		testTitles:     make(map[string]struct{}),
	}
}

func staticTitleValue(arg *ast.Node) (string, bool) {
	if arg == nil {
		return "", false
	}
	switch arg.Kind {
	case ast.KindStringLiteral:
		return arg.AsStringLiteral().Text, true
	case ast.KindNoSubstitutionTemplateLiteral:
		return arg.AsNoSubstitutionTemplateLiteral().Text, true
	default:
		return "", false
	}
}

func hasParameterizedModifier(fn *testfw.ParsedFnCall, parameterizedModifiers []string) bool {
	if fn == nil {
		return false
	}
	return slices.ContainsFunc(fn.Members, func(member string) bool {
		return slices.Contains(parameterizedModifiers, member)
	})
}

func NewRule(config Config) rule.Rule {
	return rule.Rule{
		Name: config.Name,
		Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
			contexts := []*titleLayer{newTitleLayer()}

			return rule.RuleListeners{
				ast.KindCallExpression: func(node *ast.Node) {
					fn := config.Parse(node, ctx)
					if fn == nil {
						return
					}

					cur := contexts[len(contexts)-1]
					if fn.Kind == testfw.FnKindDescribe {
						contexts = append(contexts, newTitleLayer())
					}

					if hasParameterizedModifier(fn, config.ParameterizedModifiers) {
						return
					}

					callExpr := node.AsCallExpression()
					if callExpr == nil || callExpr.Arguments == nil || len(callExpr.Arguments.Nodes) < 1 {
						return
					}
					arg0 := callExpr.Arguments.Nodes[0]
					title, ok := staticTitleValue(arg0)
					if !ok {
						return
					}

					if fn.Kind == testfw.FnKindTest {
						if _, ok := cur.testTitles[title]; ok {
							ctx.ReportNode(arg0, buildMultipleTestTitleMessage())
						}
						cur.testTitles[title] = struct{}{}
					}

					if fn.Kind != testfw.FnKindDescribe {
						return
					}
					if _, ok := cur.describeTitles[title]; ok {
						ctx.ReportNode(arg0, buildMultipleDescribeTitleMessage())
					}

					cur.describeTitles[title] = struct{}{}
				},
				rule.ListenerOnExit(ast.KindCallExpression): func(node *ast.Node) {
					fn := config.Parse(node, ctx)
					if fn == nil || fn.Kind != testfw.FnKindDescribe {
						return
					}
					if len(contexts) > 1 {
						contexts = contexts[:len(contexts)-1]
					}
				},
			}
		},
	}
}
