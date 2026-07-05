package no_focused_tests

import (
	"github.com/microsoft/typescript-go/shim/ast"
	jestUtils "github.com/web-infra-dev/rslint/internal/plugins/jest/utils"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	testfwNoFocusedTests "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/no_focused_tests"
	"github.com/web-infra-dev/rslint/internal/rule"
)

func parseJestFnCall(node *ast.Node, ctx rule.RuleContext) *testfw.ParsedFnCall {
	jestFnCall := jestUtils.ParseJestFnCall(node, ctx)
	if jestFnCall == nil {
		return nil
	}

	kind := testfw.FnKind("")
	switch jestFnCall.Kind {
	case jestUtils.JestFnTypeDescribe:
		kind = testfw.FnKindDescribe
	case jestUtils.JestFnTypeTest:
		kind = testfw.FnKindTest
	default:
		return nil
	}

	return &testfw.ParsedFnCall{
		Name:          jestFnCall.Name,
		LocalName:     jestFnCall.LocalName,
		Kind:          kind,
		Members:       jestFnCall.Members,
		MemberEntries: convertMemberEntries(jestFnCall.MemberEntries),
		Head: testfw.CallHead{
			Type: convertImportMode(jestFnCall.Head.Type),
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

func convertMemberEntries(entries []jestUtils.ParsedJestFnMemberEntry) []testfw.MemberEntry {
	out := make([]testfw.MemberEntry, len(entries))
	for i, entry := range entries {
		out[i] = testfw.MemberEntry{
			Name: entry.Name,
			Node: entry.Node,
		}
	}
	return out
}

func convertImportMode(mode jestUtils.JestImportMode) testfw.ImportMode {
	if mode == jestUtils.JEST_IMPORT_MODE {
		return testfw.ImportedMode
	}
	return testfw.GlobalMode
}

var NoFocusedTestsRule = testfwNoFocusedTests.NewRule(testfwNoFocusedTests.Config{
	Name:              "jest/no-focused-tests",
	Parse:             parseJestFnCall,
	FocusedNamePrefix: "f",
	FocusedReplacement: map[string]string{
		"fdescribe": "describe",
		"fit":       "it",
	},
})
