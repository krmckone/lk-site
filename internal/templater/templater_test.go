package templater

import (
	"fmt"
	"testing"

	"github.com/krmckone/ksite/internal/config"
)

func TestReplaceVars(t *testing.T) {
	cases := []struct {
		md        []byte
		p         config.Params
		expect    []byte
		expectErr bool
	}{
		// 0
		{
			[]byte("Hello"),
			config.Params{},
			[]byte("Hello"),
			false,
		},
		// 1
		{
			[]byte("Hello {{ .myName }}"),
			config.Params{"myName": "Kaleb"},
			[]byte("Hello Kaleb"),
			false,
		},
		// 2
		{
			[]byte("Hello {{ .myName }}"),
			config.Params{"myName": "Kaleb"},
			[]byte("Hello Kaleb"),
			false,
		},
		// 3
		{
			[]byte("Hello {{ .myName}}"),
			config.Params{"myName": "Kaleb"},
			[]byte("Hello Kaleb"),
			false,
		},
		// 4
		{
			[]byte("Hello {{ .myName		}}"),
			config.Params{"myName": "Kaleb"},
			[]byte("Hello Kaleb"),
			false,
		},
		// 5
		{
			[]byte("Hello {{ .myName }}. Welcome to {{ .projN }}"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to Testing"),
			false,
		},
		// 6
		{
			[]byte("Hello {{ .myName }}. Welcome to {{ .projN }"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to {{ .projN }"),
			true,
		},
		// 7
		{
			[]byte("Hello {{ .myName }}. Welcome to { .projN }"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to { .projN }"),
			false,
		},
		// 8
		{
			[]byte("Hello { .myName }}. Welcome to {{ .projN }"),
			config.Params{"myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello { .myName }}. Welcome to { .projN }"),
			true,
		},
		// 9
		{
			[]byte("{{ .greeting }}. It's good to see you, {{ .myName }}. Welcome to {{ .projN }}"),
			config.Params{"greeting": "Hello", "myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello. It's good to see you, Kaleb. Welcome to Testing"),
			false,
		},
		// 10
		{
			[]byte("{{ .greeting }} {{ .myName }}. Welcome to {{ .projN }}"),
			config.Params{"greeting": "Hello", "myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello Kaleb. Welcome to Testing"),
			false,
		},
		// 11
		{
			[]byte("{{ .firstPart }}{{ .lastPart }}, {{.myName}}. Welcome to {{ .projN }}"),
			config.Params{"firstPart": "Hel", "lastPart": "lo", "myName": "Kaleb", "projN": "Testing"},
			[]byte("Hello, Kaleb. Welcome to Testing"),
			false,
		},
	}
	for i, c := range cases {
		tName := fmt.Sprintf("%v: %s,%v", i, c.md, c.p)
		t.Run(tName, func(t *testing.T) {
			actualB, err := runTemplate(c.md, c.p)

			if c.expectErr && err == nil {
				t.Errorf("%s: should have had error", tName)
			} else if !c.expectErr {
				expect := string(c.expect)
				actual := string(actualB)
				if actual != expect {
					t.Errorf("Expected: %s, actual: %s", expect, actual)
				}
			}
		})
	}
}

func TestMakeTitle(t *testing.T) {
	cases := []struct {
		in, expect string
	}{
		{
			"page_0", "Page 0",
		}, {
			"page 1", "Page 1",
		}, {
			"index_page_zero", "Index Page Zero",
		},
	}
	for _, c := range cases {
		actual := makeTitle(c.in)
		if actual != c.expect {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}
