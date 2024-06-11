package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

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
	Icons  Params       `yaml:"icons"`
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

	config = ReadIcons(config)
	config.Template.Params["sheetsURL"] = config.Template.Styles.SheetURL
	config.Template.Params["currentYear"] = utils.GetCurrentYear()
	config.Template.Params["currentEasternTime"] = utils.GetCurrentEasternTime()
	return config
}

func ReadIcons(config Config) Config {
	for name, path := range config.Template.Icons {
		config.Template.Params[fmt.Sprintf("%sIcon", name)] = readIcon(path)
	}
	return config
}

func readIcon(name string) string {
	_, b, _, _ := runtime.Caller(0)
	absolutePath := filepath.Dir(b)
	return string(utils.ReadFile(fmt.Sprintf("%s/../../assets/icons/%s", absolutePath, name)))
}
