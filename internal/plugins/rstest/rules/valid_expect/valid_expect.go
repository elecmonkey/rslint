package valid_expect

import (
	rstestUtils "github.com/web-infra-dev/rslint/internal/plugins/rstest/utils"
	testfwValidExpect "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/valid_expect"
	"github.com/web-infra-dev/rslint/internal/rule"
)

func buildUnsupportedPollMemberMessage(member string) rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "unsupportedPollMember",
		Description: "expect.poll() does not support ." + member,
		Data:        map[string]string{"member": member},
	}
}

var unsupportedPollMembers = map[string]rule.RuleMessage{
	"rejects":                            buildUnsupportedPollMemberMessage("rejects"),
	"resolves":                           buildUnsupportedPollMemberMessage("resolves"),
	"throw":                              buildUnsupportedPollMemberMessage("throw"),
	"throws":                             buildUnsupportedPollMemberMessage("throws"),
	"toThrow":                            buildUnsupportedPollMemberMessage("toThrow"),
	"toThrowError":                       buildUnsupportedPollMemberMessage("toThrowError"),
	"matchSnapshot":                      buildUnsupportedPollMemberMessage("matchSnapshot"),
	"toMatchSnapshot":                    buildUnsupportedPollMemberMessage("toMatchSnapshot"),
	"toMatchInlineSnapshot":              buildUnsupportedPollMemberMessage("toMatchInlineSnapshot"),
	"toThrowErrorMatchingSnapshot":       buildUnsupportedPollMemberMessage("toThrowErrorMatchingSnapshot"),
	"toThrowErrorMatchingInlineSnapshot": buildUnsupportedPollMemberMessage("toThrowErrorMatchingInlineSnapshot"),
}

var ValidExpectRule = testfwValidExpect.NewRule(testfwValidExpect.Config{
	Name:        "rstest/valid-expect",
	ParseExpect: rstestUtils.ParseRstestExpectCallWithReason,
	UnsupportedMembers: map[string]map[string]rule.RuleMessage{
		"poll": unsupportedPollMembers,
	},
})
