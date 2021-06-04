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
	Title    string
	Content  []byte
	Template []byte
	Params   map[string]string
}

func (p *Page) String() string {
	return fmt.Sprintf(
		"Title: %s\nContent: %s\nTemplate: %s\nParams: %v",
		p.Title,
		string(p.Content),
		string(p.Template),
		p.Params,
	)
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

	pages, err := getPages("assets/pages", c.Template.Params)
	if err != nil {
		return err
	}

	for _, p := range pages {
		if err := p.exec(gm); err != nil {
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

func (p Page) exec(gm goldmark.Markdown) error {

	mdBuffer := new(bytes.Buffer)
	if err := gm.Convert(p.Content, mdBuffer); err != nil {
		return err
	}
	p.Params["main_content"] = mdBuffer.String()

	tmpl, err := template.New("template").Parse(string(p.Template))
	if err != nil {
		return err
	}

	templBuffer := new(bytes.Buffer)
	if err = tmpl.Execute(templBuffer, p.Params); err != nil {
		return err
	}

	utils.WriteFile(fmt.Sprintf("build/%s.html", p.Title), templBuffer.Bytes())
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

func getAssetsNoAbout(path string) ([]string, error) {
	assets, err := getAssets(path)
	if err != nil {
		return assets, err
	}
	result := []string{}
	for _, asset := range assets {
		if asset != "about" {
			result = append(result, asset)
		}
	}
	return result, nil
}

func runNavTemplate(md []byte, p config.Params) ([]byte, error) {
	funcs := map[string]interface{}{"getAssetsNoAbout": getAssetsNoAbout, "makeTitle": makeTitle}
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

func getPages(path string, p config.Params) ([]Page, error) {
	pages := make([]Page, 0)
	names, err := getAssets(path)
	if err != nil {
		return pages, err
	}
	// Override the base page template here
	basePage := utils.ReadFile("assets/base_page.html")
	for _, name := range names {
		md := utils.ReadFile(fmt.Sprintf("%s/%s.md", path, name))
		// Override any params here before making the page
		page, err := newPage(name, md, basePage, p)
		if err != nil {
			return pages, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func newPage(title string, content []byte, template []byte, params map[string]string) (Page, error) {
	return Page{
		title,
		content,
		template,
		params,
	}, nil
}

func makeTitle(assetName string) string {
	return strings.Title(strings.Join(strings.Split(assetName, "_"), " "))
}
