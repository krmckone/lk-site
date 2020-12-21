package preprocessor

import (
	"fmt"
	"testing"

	"github.com/krmckone/ksite/internal/config"
)

func TestReplaceVars(t *testing.T) {
	cases := []struct {
		md     []byte
		p      config.Params
		expect []byte
	}{
		{[]byte("Hello"), config.Params{}, []byte("Hello")},
		{[]byte("Hello {{ myName }}"), config.Params{"myName": "Kaleb"}, []byte("Hello Kaleb")},
		{[]byte("Hello {{myName }}"), config.Params{"myName": "Kaleb"}, []byte("Hello Kaleb")},
		{[]byte("Hello {{ myName}}"), config.Params{"myName": "Kaleb"}, []byte("Hello Kaleb")},
		{[]byte("Hello {{ myName		}}"), config.Params{"myName": "Kaleb"}, []byte("Hello Kaleb")},
		{
			[]byte("Hello {{ myName }}. Welcome to {{ projN }}"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to Testing"),
		},
		{
			[]byte("Hello {{ myName }}. Welcome to {{ projN }"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to {{ projN }"),
		},
		{
			[]byte("Hello {{ myName }}. Welcome to { projN }"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to { projN }"),
		},
		{
			[]byte("Hello { myName }}. Welcome to { projN }"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello { myName }}. Welcome to { projN }"),
		},
		{
			[]byte("{{ greeting }}. It's good to see you, {{ myName }}. Welcome to {{ projN }}"),
			config.Params{"greeting": "Hello", "myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello. It's good to see you, Kaleb. Welcome to Testing"),
		},
		{
			[]byte("{{ greeting }} {{ myName }}. Welcome to {{ projN }}"),
			config.Params{"greeting": "Hello", "myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to Testing"),
		},
		{
			[]byte("{{ firstPart }}{{ lastPart }}, {{myName}}. Welcome to {{ projN }}"),
			config.Params{"firstPart": "Hel", "lastPart": "lo", "myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello, Kaleb. Welcome to Testing"),
		},
	}
	for _, c := range cases {
		tName := fmt.Sprintf("%s,%v", c.md, c.p)
		t.Run(tName, func(t *testing.T) {
			actual := string(replaceVars(c.md, c.p))
			expect := string(c.expect)
			if actual != expect {
				t.Errorf("Expected: %s, actual: %s", expect, actual)
			}
		})
	}
}
