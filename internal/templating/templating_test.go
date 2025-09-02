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

func NewTestRuntime() utils.RuntimeConfig {
	return utils.RuntimeConfig{
		AssetsPath:  "test/assets",
		BuildPath:   "test/build",
		ConfigsPath: "test/configs",
	}
}

func TestTemplateSite(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
	if err := TemplateSite(runtime); err != nil {
		t.Errorf("Error from TemplateSite: %s", err)
	}
	_, err := os.Stat(utils.MakePath(runtime.BuildPath))
	if os.IsNotExist(err) {
		t.Errorf("%s directory not created even though TemplateSite was called: %s", runtime.BuildPath, err)
	} else if err != nil {
		t.Errorf("Error checking if %s directory exists: %s", runtime.BuildPath, err)
	}
}

func TestSetupPageParams(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := utils.Clean(utils.MakePath(runtime.BuildPath)); err != nil {
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
			[]string{filepath.Join(utils.MakePath(runtime.AssetsPath), "components", "test_component.html")},
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
		actual, err := setupPageParams(runtime, c.componentFiles, c.config, c.mainContent)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}
