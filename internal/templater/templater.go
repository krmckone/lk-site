package templater

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/krmckone/ksite/internal/config"
	"github.com/krmckone/ksite/internal/utils"
	"github.com/yuin/goldmark"
)

// BuildSite is for building the site
func BuildSite(gm goldmark.Markdown, c *config.Config) error {
	if err := runComponents(gm, c); err != nil {
		return err
	}

	if err := runPage(gm, c, "index"); err != nil {
		return err
	}

	if err := makePage(gm, c, "index"); err != nil {
		return err
	}

	return nil
}

func makePage(gm goldmark.Markdown, c *config.Config, name string) error {
	html, err := runBase(gm, c)
	if err != nil {
		return err
	}

	utils.Mkdir("build")
	utils.WriteFile(fmt.Sprintf("build/%s.html", name), html)

	return nil
}

func runComponents(gm goldmark.Markdown, c *config.Config) error {
	// Do topnav content
	if err := run(gm, c, "navbar"); err != nil {
		return err
	}

	// Do header content
	if err := run(gm, c, "header"); err != nil {
		return err
	}

	// Do footer content
	if err := run(gm, c, "footer"); err != nil {
		return err
	}

	return nil
}

func runBase(gm goldmark.Markdown, c *config.Config) ([]byte, error) {
	basePage := utils.ReadFile("assets/base_page.html")
	md, err := runTemplate(basePage, c.Template.Params)
	if err != nil {
		return md, err
	}
	return md, nil
}

func run(gm goldmark.Markdown, c *config.Config, name string) error {
	buf := new(bytes.Buffer)
	md := utils.ReadFile(fmt.Sprintf("assets/%s.md", name))
	md, err := runTemplate(md, c.Template.Params)
	if err != nil {
		return err
	}
	if err := gm.Convert(md, buf); err != nil {
		return err
	}
	c.Template.Params[name] = buf.String()
	return nil
}

func runPage(gm goldmark.Markdown, c *config.Config, name string) error {
	buf := new(bytes.Buffer)
	md := utils.ReadFile(fmt.Sprintf("assets/%s.md", name))
	md, err := runTemplate(md, c.Template.Params)
	if err != nil {
		return err
	}
	if err := gm.Convert(md, buf); err != nil {
		return err
	}
	c.Template.Params["main_content"] = buf.String()
	return nil
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
