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
