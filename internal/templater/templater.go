package templater

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/krmckone/ksite/internal/config"
)

// Run executes main template on md in environment p
func Run(md []byte, p config.Params) ([]byte, error) {
	return runTemplate(md, p)
}

// RunNav executes template for top navbar
func RunNav(md []byte, p config.Params) ([]byte, error) {
	return runNavTemplate(md, p)
}

func runTemplate(md []byte, p config.Params) ([]byte, error) {
	tmpl, err := template.New("template").Parse(string(md))
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	if err = tmpl.Execute(buffer, p); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func getAssets(path string) ([]string, error) {
	var assets []string
	dir, err := os.Open(path)
	if err != nil {
		return assets, err
	}
	assets, err = dir.Readdirnames(0)
	if err != nil {
		return assets, err
	}
	for i := range assets {
		assets[i] = strings.Split(assets[i], ".")[0]
	}
	return assets, nil
}

func runNavTemplate(md []byte, p config.Params) ([]byte, error) {
	funcs := map[string]interface{}{"getAssets": getAssets}
	tmpl, err := template.New("topnav").Funcs(funcs).Parse(string(md))
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	if err = tmpl.Execute(buffer, p); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
