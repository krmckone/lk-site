package templating

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/krmckone/lk-site/internal/utils"
)

func TestTemplateSite(t *testing.T) {
	if err := TemplateSite(); err != nil {
		t.Errorf("Error from TemplateSite: %s", err)
	}
	utils.Clean(utils.MakePath("build"))
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
			filepath.Join(utils.MakePath("assets"), "test", "pages"),
			[]string{"/pages/post_0", "/pages/post_1", "/pages/post_2"},
		},
	}
	for _, c := range cases {
		actual, err := makeHrefs(c.path)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}
