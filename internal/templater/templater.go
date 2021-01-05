package templater

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/krmckone/ksite/internal/config"
	"github.com/krmckone/ksite/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Page holds data for templating a page
type Page struct {
	Content  []byte
	Template []byte
	Params   map[string]string
}

// BuildSite is for building the site
func BuildSite() error {
	utils.Clean("build")
	utils.Mkdir("build")

	gm := newGoldmark()

	c := config.ReadConfig("configs/config.yml")

	if err := runComponents(gm, &c); err != nil {
		return err
	}

	pages, err := getAssets("assets/pages")
	if err != nil {
		return err
	}

	for _, p := range pages {
		if err := runPage(gm, &c, p); err != nil {
			return err
		}

		if err := makePage(gm, &c, p); err != nil {
			return err
		}
	}

	return nil
}

func newGoldmark() goldmark.Markdown {
	return goldmark.New(
		attributes.Enable,
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
}

func makePage(gm goldmark.Markdown, c *config.Config, name string) error {
	html, err := runBaseTemplate(gm, c)
	if err != nil {
		return err
	}

	utils.WriteFile(fmt.Sprintf("build/%s.html", name), html)

	return nil
}

func runComponents(gm goldmark.Markdown, c *config.Config) error {
	// Do topnav content
	if err := runComponentTemplate(gm, c, "topnav"); err != nil {
		return err
	}

	// Do header content
	if err := runComponentTemplate(gm, c, "header"); err != nil {
		return err
	}

	// Do footer content
	if err := runComponentTemplate(gm, c, "footer"); err != nil {
		return err
	}

	return nil
}

func runBaseTemplate(gm goldmark.Markdown, c *config.Config) ([]byte, error) {
	basePage := utils.ReadFile("assets/base_page.html")
	md, err := runTemplate(basePage, c.Template.Params)
	if err != nil {
		return md, err
	}
	return md, nil
}

func runComponentTemplate(gm goldmark.Markdown, c *config.Config, name string) error {
	buf := new(bytes.Buffer)
	md := utils.ReadFile(fmt.Sprintf("assets/%s.md", name))
	var err error
	if name == "topnav" {
		md, err = runNavTemplate(md, c.Template.Params)
	} else {
		md, err = runTemplate(md, c.Template.Params)
	}
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
	md := utils.ReadFile(fmt.Sprintf("assets/pages/%s.md", name))
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
