package valid_title

import (
	jestUtils "github.com/web-infra-dev/rslint/internal/plugins/jest/utils"
	testfwValidTitle "github.com/web-infra-dev/rslint/internal/plugins/testfw/rules/valid_title"
)

func trimFXPrefix(word string) string {
	if word == "" {
		return ""
	}
	if word[0] == 'f' || word[0] == 'x' {
		return word[1:]
	}
	return word
}

// ValidTitleRule enforces ESLint jest/valid-title.
var ValidTitleRule = testfwValidTitle.NewRule(testfwValidTitle.Config{
	Name:                   "jest/valid-title",
	Parse:                  jestUtils.ParseJestTestFnCall,
	ParameterizedModifiers: []string{"each"},
	AllowedEachSpecifiers:  "psdifjo#$%",
	TrimTitlePrefix:        trimFXPrefix,
})
