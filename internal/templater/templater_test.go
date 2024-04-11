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

func TestMakeNavTitleFromHref(t *testing.T) {
	cases := []struct {
		in, expect string
	}{
		{
			"posts/page_0", "Page 0",
		}, {
			"posts/page 1", "Page 1",
		}, {
			"index_page_zero", "Index Page Zero",
		},
	}
	for _, c := range cases {
		actual := makeNavTitleFromHref(c.in)
		if actual != c.expect {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}

func TestMakeHref(t *testing.T) {
	cases := []struct {
		assetName, originalPath, expect string
	}{
		{
			"testing_1_2_3", "test/posts", "/posts/testing_1_2_3",
		},
		{
			"file_name", "test123/xyz", "/xyz/file_name",
		},
		{
			"index", "website/test123/files", "/files/index",
		},
	}
	for _, c := range cases {
		actual := makeHref(c.assetName, c.originalPath)
		if actual != c.expect {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}

func TestMakeHrefs(t *testing.T) {
	cases := []struct {
		path   string
		expect []string
	}{
		{
			"../../assets/test/pages",
			[]string{"/pages/post_0", "/pages/post_1", "/pages/post_2"},
		},
	}
	for _, c := range cases {
		actual, err := makeHrefs(c.path)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if len(actual) != len(c.expect) {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
		for i, v := range c.expect {
			if v != actual[i] {
				t.Errorf("Expected: %s, actual: %s", c.expect, actual)
			}
		}
	}
}
