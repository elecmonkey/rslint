package utils

import (
	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	"github.com/web-infra-dev/rslint/internal/rule"
)

func ParseJestTestFnCall(node *ast.Node, ctx rule.RuleContext) *testfw.ParsedFnCall {
	jestFnCall := ParseJestFnCall(node, ctx)
	if jestFnCall == nil {
		return nil
	}

	kind := testfw.FnKind("")
	switch jestFnCall.Kind {
	case JestFnTypeDescribe:
		kind = testfw.FnKindDescribe
	case JestFnTypeTest:
		kind = testfw.FnKindTest
	default:
		return nil
	}

	return &testfw.ParsedFnCall{
		Name:          jestFnCall.Name,
		LocalName:     jestFnCall.LocalName,
		Kind:          kind,
		Members:       jestFnCall.Members,
		MemberEntries: ConvertMemberEntries(jestFnCall.MemberEntries),
		Head: testfw.CallHead{
			Type: ConvertImportMode(jestFnCall.Head.Type),
			Local: testfw.CallHeadEntry{
				Value: jestFnCall.Head.Local.Value,
				Node:  jestFnCall.Head.Local.Node,
			},
			Original: testfw.CallHeadEntry{
				Value: jestFnCall.Head.Original.Value,
				Node:  jestFnCall.Head.Original.Node,
			},
		},
	}
}

func ConvertMemberEntries(entries []ParsedJestFnMemberEntry) []testfw.MemberEntry {
	out := make([]testfw.MemberEntry, len(entries))
	for i, entry := range entries {
		out[i] = testfw.MemberEntry{
			Name: entry.Name,
			Node: entry.Node,
		}
	}
	return out
}

func ConvertImportMode(mode JestImportMode) testfw.ImportMode {
	if mode == JEST_IMPORT_MODE {
		return testfw.ImportedMode
	}
	return testfw.GlobalMode
}

func ParseJestExpectCallWithReason(node *ast.Node, ctx rule.RuleContext) (*testfw.ParsedExpectCall, string) {
	parsed := ParseJestFnCall(node, ctx)
	if parsed != nil {
		if parsed.Kind == JestFnTypeExpect {
			return ConvertExpectCall(parsed), testfw.ExpectParseReasonNone
		}
		return nil, testfw.ExpectParseReasonNone
	}

	if node == nil || node.Kind != ast.KindCallExpression {
		return nil, testfw.ExpectParseReasonNone
	}

	entries := GetJestFnMemberEntries(node)
	if len(entries) == 0 {
		return nil, testfw.ExpectParseReasonNone
	}

	if resolveExpectName(node, entries[0].Name, ctx) != "expect" {
		return nil, testfw.ExpectParseReasonNone
	}

	memberEntries := ConvertMemberEntries(entries[1:])
	_, _, reason := testfw.FindExpectModifiersAndMatcher(memberEntries)
	if reason == testfw.ExpectParseReasonMatcherNotFound && IsMemberAccessNode(node.Parent) {
		reason = testfw.ExpectParseReasonMatcherNotCalled
	}
	if reason != testfw.ExpectParseReasonNone && testfw.FindTopMostCallExpression(node) != node {
		return nil, testfw.ExpectParseReasonNone
	}

	return nil, reason
}

func ConvertExpectCall(parsed *ParsedJestFnCall) *testfw.ParsedExpectCall {
	memberEntries := ConvertMemberEntries(parsed.MemberEntries)
	out := &testfw.ParsedExpectCall{
		Head: testfw.CallHead{
			Type: ConvertImportMode(parsed.Head.Type),
			Local: testfw.CallHeadEntry{
				Value: parsed.Head.Local.Value,
				Node:  parsed.Head.Local.Node,
			},
			Original: testfw.CallHeadEntry{
				Value: parsed.Head.Original.Value,
				Node:  parsed.Head.Original.Node,
			},
		},
		Members:       parsed.Members,
		MemberEntries: memberEntries,
		ExpectCall:    parsed.Head.Local.Node.Parent,
		Factory:       "expect",
	}
	testfw.ApplyParsedExpectCall(out)
	return out
}

func resolveExpectName(node *ast.Node, localName string, ctx rule.RuleContext) string {
	name, _, _ := ResolveJestFunctionReference(node, localName, nil, ctx)
	if name == "" {
		return ""
	}
	return ApplyGlobalJestAlias(name, ctx.Settings)
}
