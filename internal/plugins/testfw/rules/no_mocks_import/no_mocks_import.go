package no_mocks_import

import (
	"fmt"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/rule"
)

const mocksDirName = "__mocks__"

type Config struct {
	Name    string
	MockAPI string
}

func buildNoManualImportErrorMessage(mockAPI string) rule.RuleMessage {
	return rule.RuleMessage{
		Id: "noManualImport",
		Description: fmt.Sprintf(
			"Mocks should not be manually imported from a %s directory. Instead use `%s` and import from the original module path",
			mocksDirName,
			mockAPI,
		),
	}
}

func isMocksImportPath(path string) bool {
	return slices.Contains(strings.Split(path, "/"), mocksDirName)
}

func isStringNode(node *ast.Node) bool {
	return node.Kind == ast.KindStringLiteral || node.Kind == ast.KindNoSubstitutionTemplateLiteral
}

func NewRule(config Config) rule.Rule {
	return rule.Rule{
		Name: config.Name,
		Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
			return rule.RuleListeners{
				ast.KindImportDeclaration: func(node *ast.Node) {
					if isMocksImportPath(node.AsImportDeclaration().ModuleSpecifier.Text()) {
						ctx.ReportNode(node, buildNoManualImportErrorMessage(config.MockAPI))
					}
				},
				ast.KindCallExpression: func(node *ast.Node) {
					callExpr := node.AsCallExpression().Expression
					if callExpr.Kind != ast.KindIdentifier {
						return
					}

					callee := callExpr.AsIdentifier()
					if callee == nil || callee.Text != "require" {
						return
					}

					arguments := node.Arguments()
					if len(arguments) == 0 {
						return
					}

					firstArg := arguments[0]
					if firstArg != nil && isStringNode(firstArg) && isMocksImportPath(firstArg.Text()) {
						ctx.ReportNode(firstArg, buildNoManualImportErrorMessage(config.MockAPI))
					}
				},
			}
		},
	}
}
