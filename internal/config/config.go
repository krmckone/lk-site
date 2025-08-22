package config

import (
	"fmt"
	"path/filepath"

	"github.com/krmckone/lk-site/internal/utils"
	"gopkg.in/yaml.v2"
)

// Config top level project config settings
type Config struct {
	Env      Params         `yaml:"environment"`
	Template TemplateConfig `yaml:"template"`
}

// TemplateConfig config for the html templating
type TemplateConfig struct {
	Params Params       `yaml:"params"`
	Icons  Params       `yaml:"icons"`
	Styles StylesParams `yaml:"styles"`
}

// Params template variable parameters config
type Params map[string]interface{}

// StylesParams parameters for stylesheets
type StylesParams struct {
	SheetURL string `yaml:"sheetURL"`
}

// ReadConfig reads in the project config yaml located at path
func ReadConfig(path string) (Config, error) {
	config := Config{}
	b, err := utils.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}

	config, err = ReadIcons(config)
	if err != nil {
		return config, err
	}
	config.Template.Params["sheetsURL"] = config.Template.Styles.SheetURL
	config.Template.Params["currentYear"] = utils.GetCurrentYear()
	config.Template.Params["currentEasternTime"] = utils.GetCurrentEasternTime()
	return config, nil
}

func ReadIcons(config Config) (Config, error) {
	for name, path := range config.Template.Icons {
		icon, err := readIcon(path)
		if err != nil {
			return config, err
		}
		config.Template.Params[fmt.Sprintf("%sIcon", name)] = icon
	}
	return config, nil
}

func readIcon(name interface{}) (string, error) {
	icon, err := utils.ReadFile(filepath.Join("assets", "icons", name.(string)))
	return string(icon), err
}
