package config

import (
	"fmt"
	"log"

	"github.com/krmckone/ksite/internal/utils"
	"gopkg.in/yaml.v2"
)

// Config top level project config settings
type Config struct {
	Template TemplateConfig `yaml:"template"`
}

// TemplateConfig config for the html templating
type TemplateConfig struct {
	Params Params       `yaml:"params"`
	Styles StylesParams `yaml:"styles"`
}

// Params template variable parameters config
type Params map[string]string

// StylesParams parameters for stylesheets
type StylesParams struct {
	SheetURL string `yaml:"sheetURL"`
}

// ReadConfig reads in the project config yaml located at path
func ReadConfig(path string) Config {
	b := utils.ReadFile(path)
	config := Config{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	config.Template.Params["githubIcon"] = ReadGitHubIcon()
	config.Template.Params["linkedinIcon"] = ReadLinkedInIcon()
	config.Template.Params["sheetsURL"] = config.Template.Styles.SheetURL
	config.Template.Params["currentYear"] = utils.GetCurrentYear()
	config.Template.Params["currentEasternTime"] = utils.GetCurrentEasternTime()
	return config
}

func ReadGitHubIcon() string {
	return readIcon("github.svg")
}

func ReadLinkedInIcon() string {
	return readIcon("linkedin.svg")
}

func readIcon(name string) string {
	return string(utils.ReadFile(fmt.Sprintf("assets/icons/%s", name)))
}
