package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/krmckone/ksite/internal/utils"
)

func TestReadConfig(t *testing.T) {
	// do setup config files
	// read them
	// test contents
	cases := []struct {
		template string
		expect   Config
	}{
		{
			`
template:
  params:
    projectName: "Hello, World!"
    myName: "Tester 0"
  styles:
    sheetURL: "styles.url"`,
			Config{
				TemplateConfig{
					Params{
						"projectName":        "Hello, World!",
						"myName":             "Tester 0",
						"sheetsURL":          "styles.url",
						"currentYear":        utils.GetCurrentYear(),
						"currentEasternTime": utils.GetCurrentEasternTime(),
					},
					StylesParams{SheetURL: "styles.url"},
				},
			},
		},
		{
			`
template:
  params:
    name: "NoName"
    yourName: "Name0"`,
			Config{
				TemplateConfig{
					Params{
						"name":               "NoName",
						"yourName":           "Name0",
						"sheetsURL":          "",
						"currentYear":        utils.GetCurrentYear(),
						"currentEasternTime": utils.GetCurrentEasternTime(),
					},
					StylesParams{},
				},
			},
		},
	}
	for _, c := range cases {
		tName := fmt.Sprintf("%v,%v", c.template, c.expect)
		t.Run(tName, func(t *testing.T) {
			utils.Mkdir("test_config")
			utils.WriteFile("test_config/config.yml", []byte(c.template))
			actual := ReadConfig("test_config/config.yml")
			if !reflect.DeepEqual(actual, c.expect) {
				t.Errorf("Expected: %v, actual: %v", c.expect, actual)
			}
			t.Cleanup(func() {
				utils.Clean("test_config")
			})
		})
	}
}
