package valid_expect

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/microsoft/typescript-go/shim/core"
	"github.com/web-infra-dev/rslint/internal/plugins/testfw"
	"github.com/web-infra-dev/rslint/internal/rule"
	internalUtils "github.com/web-infra-dev/rslint/internal/utils"
)

type validExpectOptions struct {
	AlwaysAwait   bool
	AsyncMatchers []string
	MinArgs       int
	MaxArgs       int
}

type asyncDescriptor struct {
	node           *ast.Node
	promiseWrapped bool
}

func pluralSuffix(amount int) string {
	if amount == 1 {
		return ""
	}
	return "s"
}

// Message Builders

func buildErrorTooManyArgsMessage(amount int) rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "tooManyArgs",
		Description: fmt.Sprintf("Expect takes at most %d argument%s", amount, pluralSuffix(amount)),
		Data: map[string]string{
			"amount": strconv.Itoa(amount),
			"s":      pluralSuffix(amount),
		},
	}
}

func buildErrorNotEnoughArgsMessage(amount int) rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "notEnoughArgs",
		Description: fmt.Sprintf("Expect requires at least %d argument%s", amount, pluralSuffix(amount)),
		Data: map[string]string{
			"amount": strconv.Itoa(amount),
			"s":      pluralSuffix(amount),
		},
	}
}

func buildErrorModifierUnknownMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "modifierUnknown",
		Description: "Expect has an unknown modifier",
	}
}

func buildErrorMatcherNotFoundMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "matcherNotFound",
		Description: "Expect must have a corresponding matcher call",
	}
}

func buildErrorMatcherNotCalledMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "matcherNotCalled",
		Description: "Matchers must be called to assert",
	}
}

func buildErrorAsyncMustBeAwaitedMessage(alwaysAwait bool) rule.RuleMessage {
	orReturned := " or returned"
	if alwaysAwait {
		orReturned = ""
	}
	return rule.RuleMessage{
		Id:          "asyncMustBeAwaited",
		Description: "Async assertions must be awaited" + orReturned,
		Data: map[string]string{
			"orReturned": orReturned,
		},
	}
}

func buildErrorPromisesWithAsyncAssertionsMustBeAwaitedMessage(alwaysAwait bool) rule.RuleMessage {
	orReturned := " or returned"
	if alwaysAwait {
		orReturned = ""
	}
	return rule.RuleMessage{
		Id:          "promisesWithAsyncAssertionsMustBeAwaited",
		Description: "Promises which return async assertions must be awaited" + orReturned,
		Data: map[string]string{
			"orReturned": orReturned,
		},
	}
}

func buildAsyncDescriptorMessage(descriptor asyncDescriptor, alwaysAwait bool) rule.RuleMessage {
	if descriptor.promiseWrapped {
		return buildErrorPromisesWithAsyncAssertionsMustBeAwaitedMessage(alwaysAwait)
	}
	return buildErrorAsyncMustBeAwaitedMessage(alwaysAwait)
}

func parseOptions(options any) validExpectOptions {
	out := validExpectOptions{
		AlwaysAwait:   false,
		AsyncMatchers: []string{"toReject", "toResolve"},
		MinArgs:       1,
		MaxArgs:       1,
	}

	m := internalUtils.GetOptionsMap(options)
	if m == nil {
		return out
	}

	if raw, ok := m["alwaysAwait"].(bool); ok {
		out.AlwaysAwait = raw
	}
	out.MinArgs = readIntOption(m, "minArgs", out.MinArgs)
	out.MaxArgs = readIntOption(m, "maxArgs", out.MaxArgs)
	if raw, ok := m["asyncMatchers"].([]interface{}); ok {
		out.AsyncMatchers = out.AsyncMatchers[:0]
		for _, value := range raw {
			if s, ok := value.(string); ok {
				out.AsyncMatchers = append(out.AsyncMatchers, s)
			}
		}
	}

	return out
}

func readIntOption(options map[string]interface{}, key string, defaultValue int) int {
	raw, ok := options[key]
	if !ok {
		return defaultValue
	}

	switch value := raw.(type) {
	case float64:
		return int(value)
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	default:
		return defaultValue
	}
}

type ParseExpectFn func(node *ast.Node, ctx rule.RuleContext) (*testfw.ParsedExpectCall, string)

type Config struct {
	Name               string
	ParseExpect        ParseExpectFn
	UnsupportedMembers map[string]map[string]rule.RuleMessage
}

func shouldBeAwaited(parsed *testfw.ParsedExpectCall, asyncMatchers []string) bool {
	if parsed.Async {
		return true
	}
	for _, modifier := range parsed.Modifiers {
		if modifier != "not" {
			return true
		}
	}
	return slices.Contains(asyncMatchers, parsed.Matcher)
}

func isPromiseMethodCall(node *ast.Node) bool {
	if node == nil || node.Kind != ast.KindCallExpression {
		return false
	}

	callee := ast.SkipParentheses(node.AsCallExpression().Expression)
	if !testfw.IsMemberAccessNode(callee) {
		return false
	}

	return testfw.CalleeChainName(internalUtils.AccessExpressionObject(callee)) == "Promise"
}

func getPromiseCallExpressionNode(node *ast.Node) *ast.Node {
	if node == nil {
		return nil
	}

	if node.Kind == ast.KindArrayLiteralExpression && node.Parent != nil && node.Parent.Kind == ast.KindCallExpression {
		node = node.Parent
	}

	if isPromiseMethodCall(node) {
		return node
	}

	return nil
}

func findPromiseCallExpressionNode(node *ast.Node) *ast.Node {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	if node.Parent.Kind != ast.KindCallExpression && node.Parent.Kind != ast.KindArrayLiteralExpression {
		return nil
	}
	return getPromiseCallExpressionNode(node.Parent)
}

func getParentIfPromiseChained(node *ast.Node) *ast.Node {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return node
	}

	grandParent := node.Parent.Parent
	if grandParent.Kind != ast.KindCallExpression || !testfw.IsMemberAccessNode(grandParent.AsCallExpression().Expression) {
		return node
	}

	member := grandParent.AsCallExpression().Expression
	entries := testfw.GetMemberEntries(member)
	if len(entries) == 0 {
		return node
	}

	last := entries[len(entries)-1].Name
	if last == "then" || last == "catch" {
		return getParentIfPromiseChained(grandParent)
	}

	return node
}

func isAcceptableReturnNode(node *ast.Node, allowReturn bool) bool {
	if node == nil {
		return false
	}

	if allowReturn && node.Kind == ast.KindReturnStatement {
		return true
	}
	if node.Kind == ast.KindConditionalExpression {
		return isAcceptableReturnNode(node.Parent, allowReturn)
	}

	return node.Kind == ast.KindArrowFunction || node.Kind == ast.KindAwaitExpression
}

func promiseArrayExceptionKey(sourceFile *ast.SourceFile, node *ast.Node) string {
	if sourceFile == nil || node == nil {
		return ""
	}
	r := internalUtils.TrimNodeTextRange(sourceFile, node)
	return fmt.Sprintf("%d:%d", r.Pos(), r.End())
}

func expectOpenParenRange(sourceFile *ast.SourceFile, call *ast.Node) core.TextRange {
	if sourceFile == nil || call == nil || call.Kind != ast.KindCallExpression {
		return internalUtils.TrimNodeTextRange(sourceFile, call)
	}

	callExpr := call.AsCallExpression()
	start := internalUtils.TrimNodeTextRange(sourceFile, callExpr.Expression).End()
	text := sourceFile.Text()
	for i := start; i < len(text) && i < call.End(); i++ {
		if text[i] == '(' {
			return core.NewTextRange(i, i+1)
		}
	}

	return internalUtils.TrimNodeTextRange(sourceFile, call).WithEnd(start)
}

func tooManyArgsRange(sourceFile *ast.SourceFile, args []*ast.Node, maxArgs int) core.TextRange {
	start := internalUtils.TrimNodeTextRange(sourceFile, args[maxArgs]).Pos()
	end := internalUtils.TrimNodeTextRange(sourceFile, args[len(args)-1]).End()
	if end > start {
		end--
	}
	return core.NewTextRange(start, end)
}

func asyncInsertFix(sourceFile *ast.SourceFile, fn *ast.Node) rule.RuleFix {
	switch fn.Kind {
	case ast.KindMethodDeclaration, ast.KindConstructor, ast.KindGetAccessor, ast.KindSetAccessor:
		head := internalUtils.GetFunctionHeadLoc(sourceFile, fn)
		return rule.RuleFixReplaceRange(core.NewTextRange(head.Pos(), head.Pos()), "async ")
	default:
		return rule.RuleFixInsertBefore(sourceFile, fn, "async ")
	}
}

func awaitFix(sourceFile *ast.SourceFile, node *ast.Node, alwaysAwait bool) rule.RuleFix {
	if alwaysAwait && node.Parent != nil && node.Parent.Kind == ast.KindReturnStatement {
		ret := node.Parent
		retRange := internalUtils.TrimNodeTextRange(sourceFile, ret)
		nodeRange := internalUtils.TrimNodeTextRange(sourceFile, node)
		return rule.RuleFixReplaceRange(core.NewTextRange(retRange.Pos(), nodeRange.Pos()), "await ")
	}
	return rule.RuleFixInsertBefore(sourceFile, node, "await ")
}

func resolveAsyncAssertionReportNode(
	matcherEntry *testfw.MemberEntry,
	alwaysAwait bool,
) (reportNode *ast.Node, promiseWrapped bool, insideAssertionArray bool, shouldReport bool) {
	if matcherEntry == nil || matcherEntry.Node == nil || matcherEntry.Node.Parent == nil {
		return nil, false, false, false
	}

	matcherMemberNode := matcherEntry.Node.Parent
	if matcherMemberNode.Parent == nil {
		return nil, false, false, false
	}

	promiseChainedAssertionNode := getParentIfPromiseChained(matcherMemberNode.Parent)
	insideAssertionArray = promiseChainedAssertionNode.Parent != nil && promiseChainedAssertionNode.Parent.Kind == ast.KindArrayLiteralExpression
	reportNode = promiseChainedAssertionNode
	if promiseCallNode := findPromiseCallExpressionNode(promiseChainedAssertionNode); promiseCallNode != nil {
		reportNode = promiseCallNode
		promiseWrapped = true
	}

	if reportNode.Parent == nil || isAcceptableReturnNode(reportNode.Parent, !alwaysAwait) {
		return reportNode, promiseWrapped, insideAssertionArray, false
	}

	return reportNode, promiseWrapped, insideAssertionArray, true
}

func reportAsyncDescriptor(
	ctx rule.RuleContext,
	descriptor asyncDescriptor,
	alwaysAwait bool,
	asyncInserted map[*ast.Node]bool,
) {
	msg := buildAsyncDescriptorMessage(descriptor, alwaysAwait)

	var fixes []rule.RuleFix
	if fn := ast.GetContainingFunction(descriptor.node); fn != nil {
		if !ast.IsAsyncFunction(fn) && !asyncInserted[fn] {
			fixes = append(fixes, asyncInsertFix(ctx.SourceFile, fn))
			asyncInserted[fn] = true
		}
		fixes = append(fixes, awaitFix(ctx.SourceFile, descriptor.node, alwaysAwait))
	}

	if len(fixes) > 0 {
		ctx.ReportNodeWithFixes(descriptor.node, msg, fixes...)
		return
	}
	ctx.ReportNode(descriptor.node, msg)
}

func findTopLevelMemberAccess(node *ast.Node) *ast.Node {
	current := node
	for current != nil && testfw.IsMemberAccessNode(current.Parent) {
		current = current.Parent
	}
	return current
}

func unsupportedMemberReport(parsed *testfw.ParsedExpectCall, unsupported map[string]map[string]rule.RuleMessage) (*testfw.MemberEntry, rule.RuleMessage, bool) {
	if parsed == nil || unsupported == nil {
		return nil, rule.RuleMessage{}, false
	}
	byMember := unsupported[parsed.Factory]
	if len(byMember) == 0 {
		return nil, rule.RuleMessage{}, false
	}
	for i := range parsed.MemberEntries {
		entry := &parsed.MemberEntries[i]
		if msg, ok := byMember[entry.Name]; ok {
			return entry, msg, true
		}
	}
	return nil, rule.RuleMessage{}, false
}

func NewRule(config Config) rule.Rule {
	return rule.Rule{
		Name: config.Name,
		Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
			opts := parseOptions(options)
			arrayExceptions := map[string]bool{}
			asyncInserted := map[*ast.Node]bool{}
			var descriptors []asyncDescriptor

			return rule.RuleListeners{
				ast.KindCallExpression: func(node *ast.Node) {
					parsed, reason := config.ParseExpect(node, ctx)
					if parsed == nil {
						if reason == "" {
							return
						}

						reportNode := node
						if testfw.IsMemberAccessNode(node.Parent) {
							topMember := findTopLevelMemberAccess(node.Parent)
							if topMember != nil {
								reportNode = topMember
							}
						}

						switch reason {
						case testfw.ExpectParseReasonMatcherNotFound:
							ctx.ReportNode(reportNode, buildErrorMatcherNotFoundMessage())
						case testfw.ExpectParseReasonMatcherNotCalled:
							entries := testfw.GetMemberEntries(reportNode)
							last := entries[len(entries)-1]
							if testfw.ExpectModifierNames[last.Name] {
								ctx.ReportNode(last.Node, buildErrorMatcherNotFoundMessage())
								return
							}
							ctx.ReportNode(last.Node, buildErrorMatcherNotCalledMessage())
						case testfw.ExpectParseReasonModifierUnknown:
							ctx.ReportNode(reportNode, buildErrorModifierUnknownMessage())
						}
						return
					}
					expectCall := parsed.ExpectCall
					if expectCall == nil || expectCall.Kind != ast.KindCallExpression {
						return
					}

					args := expectCall.AsCallExpression().Arguments.Nodes
					if len(args) < opts.MinArgs {
						ctx.ReportRange(
							expectOpenParenRange(ctx.SourceFile, expectCall),
							buildErrorNotEnoughArgsMessage(opts.MinArgs),
						)
					}
					if len(args) > opts.MaxArgs {
						ctx.ReportRange(
							tooManyArgsRange(ctx.SourceFile, args, opts.MaxArgs),
							buildErrorTooManyArgsMessage(opts.MaxArgs),
						)
					}

					if entry, msg, ok := unsupportedMemberReport(parsed, config.UnsupportedMembers); ok {
						ctx.ReportNode(entry.Node, msg)
						return
					}

					if parsed.MatcherEntry == nil || !shouldBeAwaited(parsed, opts.AsyncMatchers) {
						return
					}

					reportNode, promiseWrapped, insideAssertionArray, shouldReport := resolveAsyncAssertionReportNode(
						parsed.MatcherEntry,
						opts.AlwaysAwait,
					)
					if reportNode == nil {
						return
					}

					reportNodeKey := promiseArrayExceptionKey(ctx.SourceFile, reportNode)
					if shouldReport && !arrayExceptions[reportNodeKey] {
						descriptors = append(descriptors, asyncDescriptor{
							node:           reportNode,
							promiseWrapped: promiseWrapped,
						})
					}

					if insideAssertionArray {
						arrayExceptions[reportNodeKey] = true
					}
				},
				rule.ListenerOnExit(ast.KindEndOfFile): func(node *ast.Node) {
					_ = node
					for _, descriptor := range descriptors {
						reportAsyncDescriptor(ctx, descriptor, opts.AlwaysAwait, asyncInserted)
					}
				},
			}
		},
	}

}
