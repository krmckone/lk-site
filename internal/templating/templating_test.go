package templating

import (
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/krmckone/lk-site/internal/config"
	"github.com/krmckone/lk-site/internal/utils"
)

func TestTemplateSite(t *testing.T) {
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath("build")); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
	if err := TemplateSite(); err != nil {
		t.Errorf("Error from TemplateSite: %s", err)
	}
	_, err := os.Stat(utils.MakePath("build"))
	if os.IsNotExist(err) {
		t.Errorf("build directory not created even though TemplateSite was called: %s", err)
	} else if err != nil {
		t.Errorf("Error checking if build directory exists: %s", err)
	}
}

func TestMakeNavTitleFromHref(t *testing.T) {
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath("build")); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
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
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath("build")); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
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
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath("build")); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
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

func TestSetupPageParams(t *testing.T) {
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath("build")); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
	cases := []struct {
		componentFiles []string
		config         config.Config
		mainContent    string
		expect         map[string]interface{}
	}{
		{
			[]string{filepath.Join(utils.MakePath("assets"), "components", "steam_deck_top_50.html")},
			config.Config{
				Env: config.EnvConfig{Params: config.Params{}},
				Template: config.TemplateConfig{
					Params: config.Params{
						"title": "Test Page",
					},
				},
			},
			"<h1>Test Page</h1>",
			map[string]interface{}{"title": "Test Page", "main_content": template.HTML("<h1>Test Page</h1>")},
		},
	}
	for _, c := range cases {
		actual, err := setupPageParams(c.componentFiles, c.config, c.mainContent)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}
