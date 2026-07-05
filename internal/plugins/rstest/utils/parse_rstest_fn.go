package utils

import (
	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	"github.com/web-infra-dev/rslint/internal/rule"
)

var rstestMethodNames = map[string]testfw.FnKind{
	"describe": testfw.FnKindDescribe,
	"it":       testfw.FnKindTest,
	"test":     testfw.FnKindTest,
}

func ParseRstestFnCall(node *ast.Node, ctx rule.RuleContext) *testfw.ParsedFnCall {
	if node == nil || node.Kind != ast.KindCallExpression {
		return nil
	}

	memberEntries := testfw.GetMemberEntries(node)
	if len(memberEntries) == 0 {
		return nil
	}

	localName := memberEntries[0].Name
	members := make([]string, 0, len(memberEntries)-1)
	for _, entry := range memberEntries[1:] {
		members = append(members, entry.Name)
	}

	callExpr := node.AsCallExpression()
	if IsRstestFactoryCall(callExpr, members) {
		return nil
	}

	localNode := resolveHeadLocalNode(callExpr)
	name, originalNode, headType := resolveRstestFunctionReference(node, localName, localNode, ctx)
	if name == "" {
		return nil
	}

	kind, ok := rstestMethodNames[name]
	if !ok {
		return nil
	}
	if !IsValidRstestCall(name, members) {
		return nil
	}

	return &testfw.ParsedFnCall{
		Name:          name,
		LocalName:     localName,
		Kind:          kind,
		Members:       members,
		MemberEntries: memberEntries[1:],
		Head: testfw.CallHead{
			Type: headType,
			Local: testfw.CallHeadEntry{
				Value: localName,
				Node:  localNode,
			},
			Original: testfw.CallHeadEntry{
				Value: name,
				Node:  originalNode,
			},
		},
	}
}

func resolveHeadLocalNode(callExpr *ast.CallExpression) *ast.Node {
	if callExpr == nil {
		return nil
	}
	return testfw.ResolveFirstIdentifier(callExpr.Expression)
}

func resolveRstestFunctionReference(
	node *ast.Node,
	localName string,
	localNode *ast.Node,
	ctx rule.RuleContext,
) (string, *ast.Node, testfw.ImportMode) {
	if ctx.TypeChecker == nil {
		return localName, localNode, testfw.GlobalMode
	}

	callExpr := node.AsCallExpression()
	if callExpr == nil {
		return localName, localNode, testfw.GlobalMode
	}

	ident := testfw.ResolveFirstIdentifier(callExpr.Expression)
	if ident == nil || ident.Kind != ast.KindIdentifier {
		return localName, localNode, testfw.GlobalMode
	}

	symbol := ctx.TypeChecker.GetSymbolAtLocation(ident)
	if symbol == nil {
		return localName, localNode, testfw.GlobalMode
	}

	hasNonRstestImportSpecifier := false
	for _, decl := range symbol.Declarations {
		if decl == nil || decl.Kind != ast.KindImportSpecifier {
			continue
		}

		importDecl := testfw.FindImportDeclaration(decl)
		if importDecl == nil || importDecl.ModuleSpecifier == nil {
			continue
		}
		if importDecl.ModuleSpecifier.Text() != "@rstest/core" {
			hasNonRstestImportSpecifier = true
			continue
		}

		spec := decl.AsImportSpecifier()
		if spec == nil || spec.IsTypeOnly {
			continue
		}

		if spec.PropertyName != nil {
			return spec.PropertyName.Text(), spec.PropertyName, testfw.ImportedMode
		}

		name := spec.Name()
		if name != nil {
			return name.Text(), name, testfw.ImportedMode
		}
	}

	if hasNonRstestImportSpecifier {
		return "", nil, testfw.GlobalMode
	}

	return localName, localNode, testfw.GlobalMode
}

func IsRstestFactoryCall(callExpr *ast.CallExpression, members []string) bool {
	if callExpr == nil || len(members) == 0 {
		return false
	}

	last := members[len(members)-1]
	if last != "each" && last != "for" && last != "runIf" && last != "skipIf" {
		return false
	}

	switch callExpr.Expression.Kind {
	case ast.KindCallExpression, ast.KindTaggedTemplateExpression:
		return false
	default:
		return true
	}
}

func IsValidRstestCall(name string, members []string) bool {
	if _, ok := rstestMethodNames[name]; !ok {
		return false
	}

	for _, member := range members {
		switch member {
		case "only", "skip", "todo", "each", "for", "fails", "concurrent",
			"sequential", "runIf", "skipIf":
		default:
			return false
		}
	}

	return true
}
