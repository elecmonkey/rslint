package no_disabled_tests

import (
	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/plugins/jest/utils"
	testfwNoDisabledTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_disabled_tests"
	"github.com/web-infra-dev/rslint/internal/rule"
)

func isPendingCall(node *ast.Node, ctx rule.RuleContext) bool {
	if node == nil || node.Kind != ast.KindCallExpression {
		return false
	}

	callExpr := node.AsCallExpression()
	if callExpr == nil || callExpr.Expression == nil || callExpr.Expression.Kind != ast.KindIdentifier {
		return false
	}

	identifier := callExpr.Expression.AsIdentifier()
	if identifier == nil || identifier.Text != "pending" {
		return false
	}

	if ctx.TypeChecker == nil {
		return true
	}

	symbol := ctx.TypeChecker.GetSymbolAtLocation(callExpr.Expression)
	if symbol == nil {
		return true
	}

	for _, decl := range symbol.Declarations {
		if decl == nil {
			continue
		}
		if decl.Kind != ast.KindImportSpecifier {
			return false
		}

		importDecl := utils.FindImportDeclaration(decl)
		if importDecl == nil || importDecl.ModuleSpecifier == nil {
			return false
		}

		return importDecl.ModuleSpecifier.Text() == "@jest/globals"
	}

	return true
}

var NoDisabledTestsRule = testfwNoDisabledTests.NewRule(testfwNoDisabledTests.Config{
	Name:               "jest/no-disabled-tests",
	Parse:              utils.ParseJestTestFnCall,
	PendingCall:        isPendingCall,
	DisabledNamePrefix: "x",
})
