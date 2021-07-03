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
	FontURL  string `yaml:"fontURL"`
	IconURL  string `yaml:"iconURL"`
}

// ReadConfig reads in the project config yaml located at path
func ReadConfig(path string) Config {
	b := utils.ReadFile(path)
	config := Config{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	config.Template.Params["sheet_url"] = config.Template.Styles.SheetURL
	config.Template.Params["font_url"] = config.Template.Styles.FontURL
	config.Template.Params["icon_url"] = config.Template.Styles.IconURL
	fmt.Println(config.Template.Params)
	return config
}
