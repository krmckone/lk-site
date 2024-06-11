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
    sheetURL: "styles.url"
  icons:
    github: github.svg
    linkedin: linkedin.svg`,
			Config{
				TemplateConfig{
					Params{
						"projectName":        "Hello, World!",
						"myName":             "Tester 0",
						"sheetsURL":          "styles.url",
						"currentYear":        utils.GetCurrentYear(),
						"currentEasternTime": utils.GetCurrentEasternTime(),
						"githubIcon":         readIcon("github.svg"),
						"linkedinIcon":       readIcon("linkedin.svg"),
					},
					Params{
						"github":   "github.svg",
						"linkedin": "linkedin.svg",
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
					nil,
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

func TestReadIcons(t *testing.T) {
	cases := []struct {
		config Config
		expect Config
	}{
		{
			Config{
				TemplateConfig{
					Params{},
					Params{
						"github":   "github.svg",
						"linkedin": "linkedin.svg",
					},
					StylesParams{},
				},
			},
			Config{
				TemplateConfig{
					Params{
						"githubIcon":   readIcon("github.svg"),
						"linkedinIcon": readIcon("linkedin.svg"),
					},
					Params{
						"github":   "github.svg",
						"linkedin": "linkedin.svg",
					},
					StylesParams{},
				},
			},
		},
	}

	for _, c := range cases {
		actual := ReadIcons(c.config)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Expected: %v, actual: %v", c.expect, actual)
		}
	}
}
