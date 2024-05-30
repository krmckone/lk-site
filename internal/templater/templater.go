package templater

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/krmckone/ksite/internal/config"
	"github.com/krmckone/ksite/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Page holds data for templating a page
type Page struct {
	Title         string
	Content       []byte
	Template      []byte
	Params        map[string]string
	BuildPathRoot string
}

func (p *Page) String() string {
	return fmt.Sprintf(
		"Title: %s\nContent: %s\nTemplate: %s\nParams: %v\nBuildPathRoot: %s",
		p.Title,
		string(p.Content),
		string(p.Template),
		p.Params,
		p.BuildPathRoot,
	)
}

// BuildSite is for building the site. This includes templating HTML with markdown and
// putting images in the expected locations in the output
func BuildSite() error {
	utils.Clean("build")
	utils.Mkdir("build")
	utils.Mkdir("build/images")
	utils.Mkdir("build/js")
	utils.Mkdir("build/shaders")
	utils.CopyFiles("assets/images", "build/images")
	utils.CopyFiles("assets/js", "build/js")
	utils.CopyFiles("assets/shaders", "build/shaders")

	gm := newGoldmark()

	c := config.ReadConfig("configs/config.yml")

	if err := runComponents(gm, &c); err != nil {
		return err
	}

	basePages, err := getPages("assets/pages", c.Template.Params, "/")
	if err != nil {
		return err
	}

	posts, err := getPages("assets/pages/posts", c.Template.Params, "/posts")
	if err != nil {
		return err
	}

	for _, p := range append(basePages, posts...) {
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
			parser.WithAttribute(), // Lets you use {.att } syntax to add attributes to HTML output
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

	pathRoot := fmt.Sprintf("build/%s", p.BuildPathRoot)
	os.MkdirAll(pathRoot, 0700)
	utils.WriteFile(fmt.Sprintf("%s/%s.html", pathRoot, p.Title), templBuffer.Bytes())

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
	md := utils.ReadFile(fmt.Sprintf("assets/components/%s.md", name))

	var err error

	switch name {
	case "topnav":
		md, err = runNavTemplate(md, c.Template.Params)
	default:
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

	files, err := dir.Readdir(0)
	if err != nil {
		return assets, err
	}

	for i, v := range files {
		// We only want to treat files as assets here. If there's nested
		// directories containg more assets, then getAssets needs to get
		// called with that nested path to handle that case separately
		if !v.IsDir() {
			assets = append(assets, strings.Split(files[i].Name(), ".")[0])
		}
	}

	return assets, nil
}

func makeHrefs(path string) ([]string, error) {
	var hrefs []string

	assets, err := getAssets(path)
	if err != nil {
		return hrefs, err
	}

	sort.Strings(assets)
	for _, v := range assets {
		hrefs = append(hrefs, makeHref(v, path))
	}

	return hrefs, nil
}

func makeHref(assetName, originalPath string) string {
	pathSplit := strings.Split(originalPath, "/")
	hrefRoot := pathSplit[len(pathSplit)-1]

	return fmt.Sprintf("/%s/%s", hrefRoot, assetName)
}

func runNavTemplate(md []byte, p config.Params) ([]byte, error) {
	funcs := map[string]interface{}{"makeHrefs": makeHrefs, "makeNavTitle": makeNavTitleFromHref}
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

func getPages(path string, p config.Params, buildPathRoot string) ([]Page, error) {
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
		page, err := newPage(name, md, basePage, p, buildPathRoot)
		if err != nil {
			return pages, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func newPage(title string, content []byte, template []byte, params map[string]string, buildPathRoot string) (Page, error) {
	return Page{
		title,
		content,
		template,
		params,
		buildPathRoot,
	}, nil
}

func makeNavTitleFromHref(assetHref string) string {
	pathSplit := strings.Split(assetHref, "/")
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(
		strings.Join(strings.Split(pathSplit[len(pathSplit)-1], "_"), " "),
	)
}
