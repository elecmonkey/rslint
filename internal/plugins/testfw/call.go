package testfw

import (
	"strings"

	"github.com/microsoft/typescript-go/shim/ast"
)

type FnKind string

const (
	FnKindDescribe FnKind = "describe"
	FnKindTest     FnKind = "test"
)

type ImportMode string

const (
	GlobalMode   ImportMode = "global"
	ImportedMode            = "import"
)

type MemberEntry struct {
	Name string
	Node *ast.Node
}

type CallHead struct {
	Type     ImportMode
	Local    CallHeadEntry
	Original CallHeadEntry
}

type CallHeadEntry struct {
	Value string
	Node  *ast.Node
}

type ParsedFnCall struct {
	Name          string
	LocalName     string
	Kind          FnKind
	Members       []string
	MemberEntries []MemberEntry
	Head          CallHead
}

func JoinMemberEntries(entries []MemberEntry) string {
	if len(entries) == 0 {
		return ""
	}

	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = e.Name
	}

	return strings.Join(parts, ".")
}

func GetMemberEntries(node *ast.Node) []MemberEntry {
	if node == nil {
		return nil
	}

	switch node.Kind {
	case ast.KindIdentifier:
		return []MemberEntry{{
			Name: node.AsIdentifier().Text,
			Node: node,
		}}
	case ast.KindPropertyAccessExpression:
		property := node.AsPropertyAccessExpression()
		left := GetMemberEntries(property.Expression)
		nameNode := property.Name()
		if name := propertyName(nameNode); name != "" {
			return append(left, MemberEntry{
				Name: name,
				Node: nameNode,
			})
		}
		return left
	case ast.KindElementAccessExpression:
		element := node.AsElementAccessExpression()
		left := GetMemberEntries(element.Expression)
		nameNode := ast.SkipParentheses(element.ArgumentExpression)
		if name := elementAccessName(nameNode); name != "" {
			return append(left, MemberEntry{
				Name: name,
				Node: nameNode,
			})
		}
		return left
	case ast.KindCallExpression:
		return GetMemberEntries(node.AsCallExpression().Expression)
	case ast.KindTaggedTemplateExpression:
		return GetMemberEntries(node.AsTaggedTemplateExpression().Tag)
	default:
		return nil
	}
}

func propertyName(node *ast.Node) string {
	switch node.Kind {
	case ast.KindIdentifier:
		return node.AsIdentifier().Text
	case ast.KindPrivateIdentifier:
		return node.AsPrivateIdentifier().Text
	default:
		return ""
	}
}

func elementAccessName(node *ast.Node) string {
	if node == nil {
		return ""
	}

	node = ast.SkipParentheses(node)
	if node == nil {
		return ""
	}

	switch node.Kind {
	case ast.KindIdentifier:
		return node.AsIdentifier().Text
	case ast.KindStringLiteral:
		return node.AsStringLiteral().Text
	case ast.KindNoSubstitutionTemplateLiteral:
		return node.AsNoSubstitutionTemplateLiteral().Text
	default:
		return ""
	}
}

func ResolveFirstIdentifier(node *ast.Node) *ast.Node {
	if node == nil {
		return nil
	}

	switch node.Kind {
	case ast.KindIdentifier:
		return node
	case ast.KindCallExpression:
		return ResolveFirstIdentifier(node.AsCallExpression().Expression)
	case ast.KindPropertyAccessExpression:
		return ResolveFirstIdentifier(node.AsPropertyAccessExpression().Expression)
	case ast.KindElementAccessExpression:
		return ResolveFirstIdentifier(node.AsElementAccessExpression().Expression)
	case ast.KindTaggedTemplateExpression:
		return ResolveFirstIdentifier(node.AsTaggedTemplateExpression().Tag)
	default:
		return nil
	}
}

func FindImportDeclaration(node *ast.Node) *ast.ImportDeclaration {
	current := node
	for current != nil {
		switch current.Kind {
		case ast.KindImportDeclaration, ast.KindJSImportDeclaration:
			return current.AsImportDeclaration()
		}
		current = current.Parent
	}
	return nil
}

// CalleeChainName returns a dotted name for a call callee expression.
// Bracket notation contributes a segment only when the index is a supported
// static accessor name.
func CalleeChainName(expr *ast.Node) string {
	expr = ast.SkipParentheses(expr)
	if expr == nil {
		return ""
	}

	switch expr.Kind {
	case ast.KindIdentifier:
		return expr.AsIdentifier().Text
	case ast.KindPropertyAccessExpression:
		pa := expr.AsPropertyAccessExpression()
		left := CalleeChainName(pa.Expression)
		prop := pa.Name()
		if prop == nil {
			return left
		}
		pn := propertyName(prop)
		if left == "" || pn == "" {
			return left
		}
		return left + "." + pn
	case ast.KindElementAccessExpression:
		ea := expr.AsElementAccessExpression()
		left := CalleeChainName(ea.Expression)
		key := elementAccessName(ast.SkipParentheses(ea.ArgumentExpression))
		if left == "" || key == "" {
			return ""
		}
		return left + "." + key
	case ast.KindCallExpression:
		return CalleeChainName(expr.AsCallExpression().Expression)
	case ast.KindNewExpression:
		ne := expr.AsNewExpression()
		if ne == nil {
			return ""
		}
		return CalleeChainName(ne.Expression)
	case ast.KindTaggedTemplateExpression:
		return CalleeChainName(expr.AsTaggedTemplateExpression().Tag)
	default:
		return ""
	}
}
