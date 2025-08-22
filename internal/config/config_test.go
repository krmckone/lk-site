package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/krmckone/lk-site/internal/utils"
)

func TestReadConfig(t *testing.T) {
	// do setup config files
	// read them
	// test contents
	githubIcon, err := readIcon("github.svg")
	if err != nil {
		t.Errorf("Error loading test github icon: %s", err)
	}
	linkedinIcon, err := readIcon("linkedin.svg")
	if err != nil {
		t.Errorf("Error loading test github icon: %s", err)
	}
	cases := []struct {
		template string
		expect   Config
	}{
		{
			`
environment:
  params:
    steamId: "invalid_steam_id"
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
				EnvConfig{Params: Params{
					"steamId": "invalid_steam_id",
				}},
				TemplateConfig{
					Params{
						"projectName":        "Hello, World!",
						"myName":             "Tester 0",
						"sheetsURL":          "styles.url",
						"currentYear":        utils.GetCurrentYear(),
						"currentEasternTime": utils.GetCurrentEasternTime(),
						"githubIcon":         githubIcon,
						"linkedinIcon":       linkedinIcon,
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
				EnvConfig{},
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
			actual, err := ReadConfig("test_config/config.yml")
			if err != nil {
				t.Errorf("Error in reading config: %s", err)
			}
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
	githubIcon, err := readIcon("github.svg")
	if err != nil {
		t.Errorf("Error loading test github icon: %s", err)
	}
	linkedinIcon, err := readIcon("linkedin.svg")
	if err != nil {
		t.Errorf("Error loading test github icon: %s", err)
	}
	cases := []struct {
		config Config
		expect Config
	}{
		{
			Config{
				EnvConfig{Params{}},
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
				EnvConfig{Params{}},
				TemplateConfig{
					Params{
						"githubIcon":   githubIcon,
						"linkedinIcon": linkedinIcon,
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
		actual, err := ReadIcons(c.config)
		if err != nil {
			t.Errorf("Error reading icons: %s", err)
		}
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Expected: %v, actual: %v", c.expect, actual)
		}
	}
}
