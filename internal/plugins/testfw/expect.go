package testfw

import "github.com/microsoft/typescript-go/shim/ast"

const (
	ExpectParseReasonNone             = ""
	ExpectParseReasonMatcherNotFound  = "matcher-not-found"
	ExpectParseReasonModifierUnknown  = "modifier-unknown"
	ExpectParseReasonMatcherNotCalled = "matcher-not-called"
)

var ExpectModifierNames = map[string]bool{
	"not":      true,
	"rejects":  true,
	"resolves": true,
}

type ParsedExpectCall struct {
	Head            CallHead
	Members         []string
	MemberEntries   []MemberEntry
	Modifiers       []string
	ModifierEntries []MemberEntry
	Matcher         string
	MatcherEntry    *MemberEntry
	ExpectCall      *ast.Node
	Factory         string
	Async           bool
}

func IsMemberAccessNode(node *ast.Node) bool {
	if node == nil {
		return false
	}
	return node.Kind == ast.KindPropertyAccessExpression || node.Kind == ast.KindElementAccessExpression
}

func FindTopMostCallExpression(node *ast.Node) *ast.Node {
	if node == nil || node.Kind != ast.KindCallExpression {
		return node
	}

	top := node
	parent := node.Parent
	for parent != nil {
		if parent.Kind == ast.KindCallExpression {
			top = parent
			parent = parent.Parent
			continue
		}
		if parent.Kind != ast.KindPropertyAccessExpression &&
			parent.Kind != ast.KindElementAccessExpression {
			break
		}
		parent = parent.Parent
	}

	return top
}

func FindExpectModifiersAndMatcher(entries []MemberEntry) (
	[]MemberEntry,
	*MemberEntry,
	string,
) {
	if len(entries) == 0 {
		return nil, nil, ExpectParseReasonMatcherNotFound
	}

	modifiers := make([]MemberEntry, 0, len(entries))
	for _, member := range entries {
		parent := member.Node.Parent
		if parent == nil {
			return nil, nil, ExpectParseReasonModifierUnknown
		}

		grandparent := parent.Parent
		if grandparent != nil && grandparent.Kind == ast.KindCallExpression {
			return modifiers, &member, ExpectParseReasonNone
		}

		switch len(modifiers) {
		case 0:
			if !ExpectModifierNames[member.Name] {
				return nil, nil, ExpectParseReasonModifierUnknown
			}
		case 1:
			if member.Name != "not" {
				return nil, nil, ExpectParseReasonModifierUnknown
			}
			first := modifiers[0].Name
			if first != "rejects" && first != "resolves" {
				return nil, nil, ExpectParseReasonModifierUnknown
			}
		default:
			return nil, nil, ExpectParseReasonModifierUnknown
		}

		modifiers = append(modifiers, member)
	}

	return nil, nil, ExpectParseReasonMatcherNotFound
}

func ApplyParsedExpectCall(parsed *ParsedExpectCall) bool {
	modifierEntries, matcher, err := FindExpectModifiersAndMatcher(parsed.MemberEntries)
	if err != "" {
		return false
	}

	parsed.ModifierEntries = modifierEntries
	parsed.Matcher = matcher.Name
	parsed.MatcherEntry = matcher
	if len(modifierEntries) > 0 {
		parsed.Modifiers = make([]string, len(modifierEntries))
		for i, entry := range modifierEntries {
			parsed.Modifiers[i] = entry.Name
		}
	}
	return true
}
