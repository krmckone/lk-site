package preprocessor

import (
	"regexp"
	"strings"

	"github.com/krmckone/ksite/internal/config"
)

var varRegex = regexp.MustCompile("({{)(\\s*)[a-zA-Z0-9]{1,}(\\s*)(}})")

// Run executes preprocessor steps on md in environment p
func Run(md []byte, p config.Params) []byte {
	return replaceVars(md, p)
}

func replaceVars(md []byte, p config.Params) []byte {
	return varRegex.ReplaceAllFunc(md, func(b []byte) []byte {
		matchedStr := strings.TrimSpace(
			strings.TrimRight(
				strings.TrimLeft(
					string(b),
					"{{",
				), "}}",
			),
		)
		return []byte(p[matchedStr])
	})
}
