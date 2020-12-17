package preprocessor

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/krmckone/ksite/internal/config"
)

// Another equivalent version of this logic when parsing
// {{x}} vars from left to right:
// return string(b1) == "}" ||  string(b2) == "}"
// However this is incorrect when there is an error of the form
// {{x} i.e. unmatched right bracket, so I prefer the below version.
func inVar(b1, b2 byte) bool {
	if string(b1) == "}" {
		return string(b2) != "}"
	}
	return true
}

// Run executes preprocessor steps on md in environment p
func Run(md []byte, p config.Params) []byte {
	return replaceVars(md, p)
}

func replaceVars(md []byte, p config.Params) []byte {
	type loc struct {
		Start int
		End   int
	}
	mdStr := string(md)
	locs := map[string]loc{}
	for i := 0; i < len(mdStr); i++ {
		c1 := string(mdStr[i])
		c2 := ""
		if i+1 < len(mdStr) {
			c2 = string(mdStr[i+1])
		}
		if c1 == "{" && c2 == "{" {
			start := i
			i += 2
			vBytes := []byte{}
			for ; inVar(mdStr[i], mdStr[i+1]); i++ {
				if unicode.IsSpace(rune(mdStr[i])) {
					continue
				}
				vBytes = append(vBytes, mdStr[i])
			}
			v := string(vBytes)
			end := i + 2
			locs[v] = loc{start, end}
			vFmt := fmt.Sprintf("{{%s}}", v)
			mdStr = strings.Replace(mdStr, mdStr[start:end], vFmt, 1)
			mdStr = strings.Replace(mdStr, vFmt, p[v], -1)
		}
	}
	return []byte(mdStr)
}
